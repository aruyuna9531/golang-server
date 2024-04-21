package rudp

import (
	"bytes"
	"go_svr/utils/templated_list"
	"time"
)

// 参考 u35s/rudp
type message struct {
	buf  bytes.Buffer
	id   int
	tick time.Duration
}

type messageQueue struct {
	queue *templated_list.List[*message]
}

func (r *messageQueue) pop(id int) *message {
	if r.queue == nil {
		r.queue = templated_list.New[*message]()
	}
	if r.queue.Len() == 0 {
		return nil
	}
	m := r.queue.PopFront()
	if m == nil || id >= 0 && m.Value.id != id {
		return nil
	}
	//r.queue.Remove(m)
	return m.Value
}

func (r *messageQueue) push(m *message) {
	if r.queue == nil {
		r.queue = templated_list.New[*message]()
	}
	r.queue.PushBack(m)
}

type Package struct {
	Next *Package
	Bts  []byte
}

type packageBuffer struct {
	buf         bytes.Buffer
	packageList *templated_list.List[*Package]
}

func (pbf *packageBuffer) packRequest(min, max int, tag int) {
	if pbf.buf.Len()+5 > MAX_PACKAGE {
		pbf.newPackage()
	}
	pbf.buf.WriteByte(byte(tag))
	pbf.buf.WriteByte(byte((min & 0xff00) >> 8))
	pbf.buf.WriteByte(byte(min & 0xff))
	pbf.buf.WriteByte(byte((max & 0xff00) >> 8))
	pbf.buf.WriteByte(byte(max & 0xff))
}
func (pbf *packageBuffer) fillHeader(head, id int) {
	if head < 128 {
		pbf.buf.WriteByte(byte(head))
	} else {
		pbf.buf.WriteByte(byte(((head & 0x7f00) >> 8) | 0x80))
		pbf.buf.WriteByte(byte(head & 0xff))
	}
	pbf.buf.WriteByte(byte((id & 0xff00) >> 8))
	pbf.buf.WriteByte(byte(id & 0xff))
}
func (pbf *packageBuffer) packMessage(m *message) {
	if m.buf.Len()+4+pbf.buf.Len() >= MAX_PACKAGE {
		pbf.newPackage()
	}
	pbf.fillHeader(m.buf.Len()+TYPE_NORMAL, m.id)
	pbf.buf.Write(m.buf.Bytes())
}
func (pbf *packageBuffer) newPackage() {
	if pbf.buf.Len() <= 0 {
		return
	}
	p := &Package{Bts: make([]byte, pbf.buf.Len())}
	copy(p.Bts, pbf.buf.Bytes())
	pbf.buf.Reset()
	pbf.packageList.PushBack(p)
}
