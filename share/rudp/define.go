package rudp

import (
	"fmt"
	"time"
)

const (
	TYPE_PING = iota
	TYPE_EOF
	TYPE_CORRUPT
	TYPE_REQUEST
	TYPE_MISSING
	TYPE_NORMAL
)

const (
	MAX_MSG_HEAD    = 4
	GENERAL_PACKAGE = 576 - 60 - 8
	MAX_PACKAGE     = 0x7fff - TYPE_NORMAL
)

const (
	ERROR_NIL int32 = iota
	ERROR_EOF
	ERROR_REMOTE_EOF
	ERROR_CORRUPT
	ERROR_MSG_SIZE
)

// rudp
var sendTick = 10 * time.Millisecond // 每个tick的间隔时间（暂定）
var corruptTick = 5 * sendTick
var expiredTick = 5 * time.Minute
var sendDelayTick = 1 * sendTick
var missingTime = 1 * sendTick

func SetSendTick(tick time.Duration)      { sendTick = tick }
func SetCorruptTick(tick time.Duration)   { corruptTick = tick }
func SetExpiredTick(tick time.Duration)   { expiredTick = tick }
func SetSendDelayTick(tick time.Duration) { sendDelayTick = tick }
func SetMissingTime(miss time.Duration)   { missingTime = miss }

// rudp conn
var autoSend = true
var maxSendNumPerTick int = 500

func SetAutoSend(send bool)      { autoSend = send }
func SetMaxSendNumPerTick(n int) { maxSendNumPerTick = n }

func bitShow(n int) string {
	var ext string = "b"
	if n >= 1024 {
		n /= 1024
		ext = "Kb"
	}
	if n >= 1024 {
		n /= 1024
		ext = "Mb"
	}
	return fmt.Sprintf("%v %v", n, ext)
}

func FrameTime(fps int) time.Duration {
	return time.Second / time.Duration(fps)
}
