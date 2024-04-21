package lock_step

import (
	"go_svr/dependency"
	"go_svr/log"
	"go_svr/proto_codes/rpc"
	"go_svr/share/rudp"
	"google.golang.org/protobuf/proto"
	"time"
)

const FPS = 30

type frameData struct {
	idx  uint32
	cmds []*rpc.InputData
}

func newFrameData(index uint32) *frameData {
	f := &frameData{
		idx:  index,
		cmds: make([]*rpc.InputData, 0),
	}

	return f
}

type lockstep struct {
	frames     map[uint32]*frameData
	frameCount uint32
	ticker     *time.Ticker
}

func newLockstep() *lockstep {
	l := &lockstep{
		frames: make(map[uint32]*frameData),
		ticker: time.NewTicker(time.Second / FPS),
	}

	return l
}

var ls = newLockstep() // 实际上不会是单例 这个测试用

func GetInst() *lockstep {
	return ls
}

func (l *lockstep) Init() {
	dependency.PushCmd = ls.pushCmd
	go ls.loop()
}

func (l *lockstep) reset() {
	l.frames = make(map[uint32]*frameData)
	l.frameCount = 0
}

func (l *lockstep) getFrameCount() uint32 {
	return l.frameCount
}

func (l *lockstep) pushCmd(cmd *rpc.InputData) bool {
	f, ok := l.frames[l.frameCount]
	if !ok {
		f = newFrameData(l.frameCount)
		l.frames[l.frameCount] = f
	}
	// 检查是否同一帧发来两次操作
	for _, v := range f.cmds {
		if v.Id == cmd.Id {
			return false
		}
	}
	f.cmds = append(f.cmds, cmd)
	log.Debug("cmd pushed to frame %d, id %d, sid %d, x %d, y%d, seat %d", l.frameCount, cmd.Id, cmd.Sid, cmd.X, cmd.Y, cmd.Roomseatid)
	return true
}

func (l *lockstep) getRangeFrames(from, to uint32) []*frameData {
	ret := make([]*frameData, 0, to-from)

	for ; from <= to && from <= l.frameCount; from++ {
		f, ok := l.frames[from]
		if !ok {
			continue
		}
		ret = append(ret, f)
	}

	return ret
}

func (l *lockstep) getFrame(idx uint32) *frameData {
	return l.frames[idx]
}

func (l *lockstep) loop() {
	for {
		select {
		case <-l.ticker.C:
			f := &rpc.FrameData{
				FrameID: l.frameCount,
			}
			if fd, ok := l.frames[l.frameCount]; ok {
				f.Input = fd.cmds
			} else {
				f.Input = []*rpc.InputData{}
			}
			b, err := proto.Marshal(f)
			if err != nil {
				log.Error("lockstep loop tick error: %s", err.Error())
				continue
			}
			rudp.GetInst().Broadcast(b)
			if len(f.Input) > 0 {
				log.Debug("broadcast frame id %d to players, len = %d", l.frameCount, len(f.Input))
			}
			l.frameCount++
		}
	}
}
