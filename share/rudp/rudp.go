package rudp

import (
	"go_svr/log"
	"go_svr/utils/templated_list"
	"sync/atomic"
	"time"
)

type Rudp struct {
	recvQueue    messageQueue // 接收队列
	recvSkip     map[int]time.Duration
	reqSendAgain chan [2]int
	recvIDMin    int
	recvIDMax    int

	sendQueue    messageQueue
	sendHistory  messageQueue
	addSendAgain chan [2]int
	sendID       int64

	corrupt Error // 崩溃报告

	currentTick       time.Duration
	lastRecvTick      time.Duration
	lastExpiredTick   time.Duration
	lastSendDelayTick time.Duration
}

func New() *Rudp {
	return &Rudp{
		recvQueue:    messageQueue{queue: templated_list.New[*message]()},
		sendQueue:    messageQueue{queue: templated_list.New[*message]()},
		sendHistory:  messageQueue{queue: templated_list.New[*message]()},
		reqSendAgain: make(chan [2]int, 1<<10),
		addSendAgain: make(chan [2]int, 1<<10),
		recvSkip:     make(map[int]time.Duration),
	}
}

func (r *Rudp) Read(bts []byte) (int, error) {
	if err := r.corrupt.Load(); err != ERROR_NIL {
		log.Debug("Rudp::Read corrupted, err = %d", err)
		return 0, r.corrupt.Error()
	} // 已崩，不做事
	m := r.recvQueue.pop(r.recvIDMin) // 从接收队列弹出第1个元素，并且它的id要等于recvIDMin（这个是否有问题？底层是list的情况下，并不一定是按id顺序排序的）
	if m == nil {
		return 0, nil
	} // 当队列里没有元素或者id不为recvIDMin时返回nil
	r.recvIDMin++
	copy(bts, m.buf.Bytes()) // 把buf的内容拷贝到入参
	log.Debug("Rudp::Read copied %d bytes to bts", m.buf.Len())
	return m.buf.Len(), nil
}

func (r *Rudp) Write(bts []byte) (n int, err error) {
	if err := r.corrupt.Load(); err != ERROR_NIL {
		log.Debug("Rudp::Write 坏了，原因是%d", err)
		return 0, r.corrupt.Error()
	} // 崩溃检测，略
	if len(bts) >= MAX_PACKAGE {
		log.Debug("Rudp::Write 退出，因为入参长度(%d)超过限制(%d)", len(bts), MAX_PACKAGE)
		return 0, nil
	} // 入参已满
	m := &message{}
	n, err = m.buf.Write(bts)
	if err != nil {
		log.Debug("Rudp::Write buf写出错: %s", err.Error())
		return 0, err
	}
	m.id = int(atomic.AddInt64(&r.sendID, 1)) - 1
	m.tick = r.currentTick
	r.sendQueue.push(m)
	log.Debug("Rudp::Write 往发送队列写入了%v", m)
	return n, nil
}

// Update 更新时间片
func (r *Rudp) Update(tick time.Duration) *Package {
	if r.corrupt.Load() != ERROR_NIL {
		return nil
	}
	r.currentTick += tick // 更新 currentTick
	log.Debug("Rudp::Update 将当前tick修改为%d", r.currentTick)
	if r.currentTick >= r.lastExpiredTick+expiredTick {
		// 当前tick距离上次过期tick已经经过expiredTick周期，清理一次sendHistory的过期帧
		r.lastExpiredTick = r.currentTick
		log.Debug("Rudp::Update 将最近过期Tick更改为%d", r.lastExpiredTick)
		r.clearSendExpired()
	}
	if r.currentTick >= r.lastRecvTick+corruptTick {
		log.Debug("Rudp::Update 毁坏，原因ERROR_CORRUPT")
		r.corrupt.Store(ERROR_CORRUPT)
	}
	if r.currentTick >= r.lastSendDelayTick+sendDelayTick {
		r.lastSendDelayTick = r.currentTick
		log.Debug("Rudp::Update 将上次发送延迟tick更改为%d", r.lastSendDelayTick)
		return r.output()
	}
	return nil
}

func (r *Rudp) getID(max int, bt1, bt2 byte) int {
	n1, n2 := int(bt1), int(bt2)
	id := n1*256 + n2
	id |= max & ^0xffff
	if id < max-0x8000 {
		id += 0x10000
		log.Debug("id < max-0x8000 ,net %v,id %v,min %v,max %v,cur %v",
			n1*256+n2, id, r.recvIDMin, max, id+0x10000)
	} else if id > max+0x8000 {
		id -= 0x10000
		log.Debug("id > max-0x8000 ,net %v,id %v,min %v,max %v,cur %v",
			n1*256+n2, id, r.recvIDMin, max, id+0x10000)
	}
	log.Debug("Rudp::getID returns %d", id)
	return id
}

func (r *Rudp) output() *Package {
	var tmp = packageBuffer{
		packageList: templated_list.New[*Package](),
	}
	log.Debug("Rudp::outPut 检查重传请求")
	r.reqMissing(&tmp)
	log.Debug("Rudp::outPut 响应请求")
	r.replyRequest(&tmp)
	log.Debug("Rudp::outPut 发送消息")
	r.sendMessage(&tmp)
	if tmp.packageList.Front() == nil && tmp.buf.Len() == 0 {
		log.Debug("Rudp::outPut 发送心跳")
		tmp.buf.WriteByte(byte(TYPE_PING))
	}
	tmp.newPackage()
	return tmp.packageList.Front().Value
}

func (r *Rudp) Input(bts []byte) {
	sz := len(bts)
	if sz > 0 {
		r.lastRecvTick = r.currentTick
		log.Debug("Rudp::outPut 设置上次接收时间为%d", r.lastRecvTick)
	}
	for sz > 0 {
		btsLen := int(bts[0])
		if btsLen > 127 {
			if sz <= 1 {
				r.corrupt.Store(ERROR_MSG_SIZE)
				log.Debug("Rudp::outPut 毁坏，原因ERROR_MSG_SIZE")
				return
			}
			btsLen = (btsLen*256 + int(bts[1])) & 0x7fff
			bts = bts[2:]
			sz -= 2
		} else {
			bts = bts[1:]
			sz -= 1
		}
		switch btsLen {
		case TYPE_PING:
			log.Debug("Rudp::outPut btsLen = TYPE_PING, checkMissing")
			r.checkMissing(false)
		case TYPE_EOF:
			log.Debug("Rudp::outPut btsLen = TYPE_EOF, 毁坏，原因ERROR_EOF")
			r.corrupt.Store(ERROR_EOF)
		case TYPE_CORRUPT:
			log.Debug("Rudp::outPut btsLen = TYPE_CORRUPT, 毁坏，原因ERROR_REMOTE_EOF")
			r.corrupt.Store(ERROR_REMOTE_EOF)
			return
		case TYPE_REQUEST, TYPE_MISSING:
			log.Debug("Rudp::outPut btsLen = %d")
			if sz < 4 {
				log.Debug("Rudp::outPut sz = %d, 毁坏，原因ERROR_MSG_SIZE")
				r.corrupt.Store(ERROR_MSG_SIZE)
				return
			}
			exe := r.addRequest
			max := int(r.sendID)
			if btsLen == TYPE_MISSING {
				log.Debug("Rudp::outPut btsLen == TYPE_MISSING, 更改exe为addMissing，最大值更改为%d", r.recvIDMax)
				exe = r.addMissing
				max = r.recvIDMax
			}
			exe(r.getID(max, bts[0], bts[1]), r.getID(max, bts[2], bts[3]))
			bts = bts[4:]
			sz -= 4
		default:
			btsLen -= TYPE_NORMAL
			if sz < btsLen+2 {
				log.Debug("Rudp::outPut sz(%d) < btsLen+2(%d), 毁坏，原因ERROR_MSG_SIZE", sz, btsLen+2)
				r.corrupt.Store(ERROR_MSG_SIZE)
				return
			}
			log.Debug("Rudp::outPut 正在插入信息，bts = %s", bts)
			r.insertMessage(r.getID(r.recvIDMax, bts[0], bts[1]), bts[2:btsLen+2])
			bts = bts[btsLen+2:]
			sz -= btsLen + 2
		}
	}
	log.Debug("Rudp::outPut 检查丢失数据(false)")
	r.checkMissing(false)
}

func (r *Rudp) checkMissing(direct bool) {
	head := r.recvQueue.queue.Front()
	if head != nil && head.Value.id > r.recvIDMin {
		nano := time.Duration(time.Now().UnixNano())
		last := r.recvSkip[r.recvIDMin]
		if !direct && last == 0 {
			r.recvSkip[r.recvIDMin] = nano
			log.Debug("丢失起始数据 %v-%v,最大值 %v", r.recvIDMin, head.Value.id-1, r.recvIDMax)
		} else if direct || last+missingTime < nano {
			delete(r.recvSkip, r.recvIDMin)
			r.reqSendAgain <- [2]int{r.recvIDMin, head.Value.id - 1}
			log.Debug("需要已丢失数据 %v-%v,direct %v,等待数量 %v",
				r.recvIDMin, head.Value.id-1, direct, r.recvQueue.queue.Len())
		}
	}
}

func (r *Rudp) insertMessage(id int, bts []byte) {
	if id < r.recvIDMin {
		log.Debug("already recv %v,len %v", id, len(bts))
		return
	}
	delete(r.recvSkip, id)
	if id > r.recvIDMax || r.recvQueue.queue.Len() == 0 {
		m := &message{}
		m.buf.Write(bts)
		m.id = id
		r.recvQueue.push(m)
		r.recvIDMax = id
	} else {
		last := r.recvQueue.queue.Front()
		for m := r.recvQueue.queue.Front(); m != nil; m = m.Next() {
			if m.Value.id == id {
				// 已收到的消息，丢弃并返回（可能是对端已发但很迟才收到的帧，此前已经请求重传并收到了重传帧）
				log.Debug("repeat recv id %v,len %v", id, len(bts))
				return
			}
			if m.Value.id > id {
				tmp := &message{}
				tmp.buf.Write(bts)
				tmp.id = id
				last.Value = tmp
				return
			}
			last = m.Next()
		}
	}
}

func (r *Rudp) sendMessage(tmp *packageBuffer) {
	log.Debug("Rudp::sendMessage 发送全部消息")
	for m := r.sendQueue.queue.Front(); m != nil; m = m.Next() {
		tmp.packMessage(m.Value)
	}
	if r.sendQueue.queue.Len() > 0 {
		log.Debug("Rudp::sendMessage 将当前队列加入历史队列")
		if r.sendHistory.queue.Len() == 0 {
			r.sendHistory = r.sendQueue
		} else {
			r.sendHistory.queue.PushBackList(r.sendQueue.queue)
		}
		log.Debug("Rudp::sendMessage 发送队列清空")
		r.sendQueue.queue.Init()
	}
}

func (r *Rudp) clearSendExpired() {
	m := r.sendHistory.queue.Front()
	for m != nil {
		if m.Value.tick >= r.lastExpiredTick {
			break
		}
		oldm := m
		m = m.Next()
		r.sendHistory.queue.Remove(oldm)
	}
}

func (r *Rudp) addRequest(min, max int) {
	log.Debug("add request %v-%v,max send id %v", min, max, r.sendID)
	r.addSendAgain <- [2]int{min, max}
}

func (r *Rudp) addMissing(min, max int) {
	if max < r.recvIDMin {
		log.Debug("add missing %v-%v fail,already recv,min %v", min, max, r.recvIDMin)
		return
	}
	if min > r.recvIDMin {
		log.Debug("add missing %v-%v fail, more than min %v", min, max, r.recvIDMin)
		return
	}
	head := 0
	if r.recvQueue.queue.Len() != 0 {
		head = r.recvQueue.queue.Front().Value.id
	}
	log.Debug("add missing %v-%v,min %v,head %v", min, max, r.recvIDMin, head)
	r.recvIDMin = max + 1
	r.checkMissing(true)
}

func (r *Rudp) replyRequest(tmp *packageBuffer) {
	for {
		select {
		case again := <-r.addSendAgain:
			history := r.sendHistory.queue.Front()
			min, max := again[0], again[1]
			if history == nil || max < history.Value.id {
				log.Debug("重传，丢失了 %v-%v,send max %v", min, max, r.sendID)
				tmp.packRequest(min, max, TYPE_MISSING)
			} else {
				var start, end, num int
				for {
					if history == nil || max < history.Value.id {
						// 已过期
						break
					} else if min <= history.Value.id {
						tmp.packMessage(history.Value)
						if start == 0 {
							start = history.Value.id
						}
						end = history.Value.id
						num++
					}
					history = history.Next()
				}
				if min < start {
					tmp.packRequest(min, start-1, TYPE_MISSING)
					log.Debug("重传丢失的数据包 %v-%v,send max %v", min, start-1, r.sendID)
				}
				log.Debug("重传 %v-%v of %v-%v, 数量 %v,最大发送id %v", start, end, min, max, num, r.sendID)
			}
		default:
			return
		}
	}
}

func (r *Rudp) reqMissing(tmp *packageBuffer) {
	for {
		select {
		case req := <-r.reqSendAgain:
			log.Debug("Rudp::reqMissing 检测到重传请求，从%d~%d", req[0], req[1])
			tmp.packRequest(req[0], req[1], TYPE_REQUEST)
		default:
			return
		}
	}
}
