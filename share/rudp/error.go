package rudp

import (
	"errors"
	"sync/atomic"
)

type Error struct {
	v int32
}

func (e *Error) Load() int32   { return atomic.LoadInt32(&e.v) }
func (e *Error) Store(n int32) { atomic.StoreInt32(&e.v, n) }

func (e *Error) Error() error {
	switch e.Load() {
	case ERROR_EOF:
		return errors.New("EOF")
	case ERROR_REMOTE_EOF:
		return errors.New("remote EOF")
	case ERROR_CORRUPT:
		return errors.New("corrupt")
	case ERROR_MSG_SIZE:
		return errors.New("recive msg size error")
	default:
		return nil
	}
}
