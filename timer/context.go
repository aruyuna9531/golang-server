package timer

import (
	"context"
	"errors"
	"go_svr/log"
	"sync/atomic"
	"time"
)

var nextTimerIdForContext atomic.Int32

const InfiniteLoop = -1

type Callback func()

type ContextTrigger struct {
	id          int32
	ctx         context.Context
	cancelFunc  context.CancelFunc
	repeatTimes int32

	totalRepeat int32    // 循环总次数，执行完该次数本计时器自动删除，为-1则无限制
	repeatMilli int64    // 循环间隔（ms），不能为0
	callback    Callback // 到点执行的函数
}

func NewContextTrigger(repeatMilli int64, repeatTimes int32, callback Callback) *ContextTrigger {
	return &ContextTrigger{
		repeatMilli: repeatMilli,
		totalRepeat: repeatTimes,
		callback:    callback,
	}
}

func (ct *ContextTrigger) GetId() int32 {
	return ct.id
}

// start 启动计时器，第一次触发就是RepeatMilli指定的时间
func (ct *ContextTrigger) start() {
	ct.id = nextTimerIdForContext.Add(1)
	ct.ctx, ct.cancelFunc = context.WithTimeout(context.Background(), time.Duration(ct.repeatMilli)*time.Millisecond)
	go ct.wait()
}

// startForDelay 第一次触发是输入的delayMilli，后面以RepeatMilli循环
func (ct *ContextTrigger) startForDelay(delayMilli int64) {
	ct.id = nextTimerIdForContext.Add(1)
	ct.ctx, ct.cancelFunc = context.WithTimeout(context.Background(), time.Duration(delayMilli)*time.Millisecond)
	go ct.wait()
}

// cancel 手动取消计时器
func (ct *ContextTrigger) cancel() {
	if ct.cancelFunc != nil {
		ct.cancelFunc()
	}
}

func (ct *ContextTrigger) wait() {
	for {
		<-ct.ctx.Done()
		if errors.Is(ct.ctx.Err(), context.Canceled) {
			log.Debug("context timer id %d canceled manually, exit", ct.id)
			return
		}
		if ct.callback != nil {
			ct.callback()
		}
		ct.repeatTimes++
		if ct.totalRepeat != InfiniteLoop && ct.repeatTimes >= ct.totalRepeat {
			log.Debug("context timer id %d times is up, count %d", ct.id, ct.repeatTimes)
			ctm.RemoveTimer(ct.id)
			return
		}
		if ct.repeatMilli > 0 {
			ct.ctx, ct.cancelFunc = context.WithTimeout(context.Background(), time.Duration(ct.repeatMilli)*time.Millisecond)
			ddl, _ := ct.ctx.Deadline()
			log.Debug("context timer id %d start next loop, done at %s", ct.id, ddl.Format("2006-01-02 15:04:05"))
		} else {
			ct.cancel()
			ctm.RemoveTimer(ct.id)
		}
		log.Debug("Context timer id %d done. times %d", ct.id, ct.repeatTimes)
	}
}

type ContextTimer struct {
	tmMap map[int32]*ContextTrigger
}

var ctm = &ContextTimer{tmMap: make(map[int32]*ContextTrigger)}

func GetContextTimer() *ContextTimer {
	return ctm
}

// AddTimer 添加倒计时
// delayMilli 第一次执行的延后时间（ms） 填0是生成计时器后会立即执行一次
// repeat 重复执行的时间间隔
// times 执行限定次数（如无限制，填入timer.InfiniteLoop或-1） 0是不执行
// callback 到点后执行的函数
// return id 本次生成的计时器id（手动取消时可以作为参数）
func (tm *ContextTimer) AddTimer(delayMilli int64, repeat int64, times int32, callback Callback) (id int32) {
	if times == 0 {
		return 0
	}
	if repeat <= 0 {
		repeat = 1
	}
	newT := NewContextTrigger(repeat, times, callback)
	newT.startForDelay(delayMilli)
	tm.tmMap[newT.GetId()] = newT
	log.Debug("timer %d added", newT.GetId())
	return newT.GetId()
}

// RemoveTimer 取消一个计时器
func (tm *ContextTimer) RemoveTimer(id int32) {
	ct, ok := tm.tmMap[id]
	if !ok {
		return
	}
	ct.cancel()
	delete(tm.tmMap, id)
	log.Debug("timer %d removed", id)
}
