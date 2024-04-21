package timer

import (
	"fmt"
	"github.com/aruyuna9531/skiplist"
	"go_svr/framebase"
	"go_svr/log"
	"go_svr/proto_codes/rpc"
	"go_svr/utils/sorted_set"
	"sync/atomic"
	"time"
)

var nextTimerId atomic.Int32

type Trigger struct {
	timerId    int32
	Fun        func(int64, int32, ...interface{})
	Param      []interface{}
	RepeatTime time.Duration

	repeatCounter int32
	nextTrigger   int64
}

func (t *Trigger) Run() {
	t.Fun(t.nextTrigger, t.repeatCounter, t.Param...)
	t.repeatCounter++
	if t.RepeatTime != 0 {
		t.nextTrigger += int64(t.RepeatTime / time.Millisecond)
	}
}

func (t *Trigger) Key() int32 {
	return t.timerId
}

func (t *Trigger) Less(i skiplist.ISkiplistElement[int32]) bool {
	ii, ok := i.(*Trigger)
	if !ok {
		return false
	}
	if t.nextTrigger < ii.nextTrigger {
		return true
	}
	if t.nextTrigger > ii.nextTrigger {
		return false
	}
	return t.timerId < ii.timerId
}

type Timer struct {
	triggers *sorted_set.SortedSet[int32, *Trigger]
}

func (t *Timer) PushNewTimerTrigger(firstTriggerTime int64, repeatMilli int64, cb func(int64, int32, ...interface{}), param ...interface{}) (id int32) {
	now := time.Now().UnixMilli()
	if now >= firstTriggerTime {
		cb(now, 0, param...)
		if repeatMilli == 0 {
			return 0
		}
		firstTriggerTime = now + repeatMilli
	}

	if t.triggers == nil {
		t.triggers = sorted_set.NewSortedSet[int32, *Trigger]()
	}
	newTimerId := nextTimerId.Add(1)
	_ = t.triggers.Add(&Trigger{
		timerId:       newTimerId,
		Fun:           cb,
		Param:         param,
		RepeatTime:    time.Duration(repeatMilli) * time.Millisecond,
		repeatCounter: 0,
		nextTrigger:   firstTriggerTime,
	})
	return newTimerId
}

func (t *Timer) RemoveTimer(id int32) error {
	return t.triggers.Delete(id)
}

func (t *Timer) Trigger(nowTs int64) {
	for t.triggers.GetCount() > 0 {
		ts, err := t.triggers.GetByRank(1)
		if err != nil {
			log.Error("Timer trigger error: %s", err.Error())
			return
		}
		if nowTs < ts.nextTrigger {
			log.Debug("next trigger time %d is not expired", ts.nextTrigger)
			break
		}
		ts.Run()
		if ts.RepeatTime != 0 {
			_, err := t.triggers.Update(ts)
			if err != nil {
				log.Error("update repeated timer error: %s", err.Error())
			}
		} else {
			err := t.triggers.Delete(ts.timerId)
			if err != nil {
				log.Error("Remove once timer error: %s", err.Error())
			}
		}
	}
}

var tm = &Timer{
	triggers: sorted_set.NewSortedSet[int32, *Trigger](),
}

func GetInst() *Timer {
	return tm
}

func TimerTestCode() {
	id := GetInst().PushNewTimerTrigger(0, 20000, func(now int64, repeatCount int32, _ ...interface{}) {
		// 服务器每启动20秒给所有连接的客户端一个推送
		framebase.GetTcpSvr().Broadcast(rpc.MessageId_Msg_SC_Message, &rpc.SC_Message{
			ErrCode: 0,
			Message: fmt.Sprintf("程序已启动%d秒", (repeatCount+1)*20),
		})
	})
	go func(idd int32) {
		time.Sleep(50 * time.Second)
		err := GetInst().RemoveTimer(idd)
		if err != nil {
			log.Error("Remove test timer failed: %s", err.Error())
			return
		}
		log.Info("Remove test timer success")
	}(id)
}
