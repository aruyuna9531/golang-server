package timer

import (
	"go_svr/log"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	log.Debug("Test Start")
	id := GetContextTimer().AddTimer(0, 0, InfiniteLoop, func() {
		log.Debug("test timer called")
	})
	GetContextTimer().AddTimer(2000, 5000, 3, func() {
		log.Debug("test sub timer called")
	})
	go func() {
		time.Sleep(2 * time.Second)
		GetContextTimer().RemoveTimer(id)
	}()
	time.Sleep(30 * time.Second)
}
