package util

import (
	"fmt"
	"hash"
	"time"

	capnp "github.com/glycerine/go-capnproto"
)

func NewTimeNow(n *capnp.Segment) Time {
	t := NewTime(n)
	t.SetTime(time.Now())
	return t
}

func NewTimeFrom(n *capnp.Segment, o time.Time) Time {
	t := NewTime(n)
	t.SetTime(o)
	return t
}

func (t Time) Time() time.Time {
	return time.Unix(int64(t.Seconds()), int64(t.NanoSeconds()))
}

func (t Time) SetTime(o time.Time) {
	o = o.UTC()
	t.SetSeconds(uint64(o.Unix()))
	t.SetNanoSeconds(uint32(o.Nanosecond()))
}

func (t Time) Hash(h hash.Hash) {
	h.Write(UInt64ToBA(t.Seconds()))
	h.Write(UInt32ToBA(t.NanoSeconds()))
}

func (t Time) String() string {
	return fmt.Sprintf("Time{%v}", t.Time())
}
