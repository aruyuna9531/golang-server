package timer

import (
	"fmt"
	"go_svr/log"
	"go_svr/mytcp"
	"time"
)

type Trigger struct {
	Fun        func(int64, interface{}, int32)
	Param      interface{}
	RepeatTime time.Duration

	repeatCounter int32
}

type Timer struct {
	triggers map[int64][]Trigger //TODO ←这里实际上用的是有序列表，有时间再手撸
}

func (t *Timer) PushTimerTrigger(atTs int64, trigger Trigger) { // TODO at应为时间戳 跟上面的todo一起做
	if t.triggers == nil {
		t.triggers = make(map[int64][]Trigger)
	}
	t.triggers[atTs] = append(t.triggers[atTs], trigger)
}

func (t *Timer) Trigger(nowTs int64) {
	ts, ok := t.triggers[nowTs]
	if !ok {
		return
	}
	for _, trigger := range ts {
		trigger.Fun(time.Now().Unix(), trigger.Param, trigger.repeatCounter)
		trigger.repeatCounter++
		if trigger.RepeatTime != 0 {
			t.PushTimerTrigger(time.Now().Add(trigger.RepeatTime).Unix(), trigger)
		}
	}
	if len(ts) > 0 {
		log.Debug("timer trigger finished, event count: %d", len(ts))
	}
	delete(t.triggers, nowTs)
}

var tm = &Timer{
	triggers: map[int64][]Trigger{},
}

func GetInst() *Timer {
	return tm
}

func PushTrigger(atTs int64, trigger Trigger) {
	tm.PushTimerTrigger(atTs, trigger)
}

// 指定时间之后执行
func PushTriggerAfterDelay(delay time.Duration, trigger Trigger) {
	tm.PushTimerTrigger(time.Now().Add(delay).Unix(), trigger)
}

func TimerTestCode() {
	PushTrigger(time.Now().Add(20*time.Second).Unix(), Trigger{
		Fun: func(now int64, _ interface{}, repeatCount int32) {
			// 服务器每启动20秒给所有连接的客户端一个推送
			mytcp.GetTcpSvr().Broadcast([]byte(fmt.Sprintf("程序已启动%d秒", (repeatCount+1)*20)))
		},
		RepeatTime: 20 * time.Second,
	})
}
