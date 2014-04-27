package message

import (
	"hash"

	"util"
)

func (m MessagePayload) Hash(h hash.Hash) {
	switch m.Which() {
	case MESSAGEPAYLOAD_VOTE:
		m.Vote().Hash(h)
	case MESSAGEPAYLOAD_CHANGE:
		m.Change().Hash(h)
	default:
	}
}

func (c Change_List) Hash(h hash.Hash) {
	for i := 0; i < c.Len(); i++ {
		c.At(i).Hash(h)
	}
}

func (v Vote) Hash(h hash.Hash) {
	v.Votes().Hash(h)
	v.Rate().Hash(h)
	v.Drop().Hash(h)
	v.Time().Hash(h)
	v.Authorization().Hash(h)
}

func (c Change) Hash(h hash.Hash) {
	c.Authorization().Hash(h)
	c.Created().Hash(h)
	switch c.Type().Which() {
	case CHANGETYPE_TRANSACTION:
		c.Type().Transaction().Hash(h)
	case CHANGETYPE_LOCATION:
		c.Type().Location().Hash(h)
	default:
	}
}

func (t TransactionChange) Hash(h hash.Hash) {
	h.Write(t.Source())
	h.Write(t.Destination())
	h.Write(util.UInt64ToBA(t.Amount()))
}

func (l LocationChange) Hash(h hash.Hash) {
	h.Write(l.Account())
	h.Write([]byte(l.Location()))
}

func (dl DropChange_List) Hash(h hash.Hash) {
	for i := 0; i < dl.Len(); i++ {
		dl.At(i).Hash(h)
	}
}

func (d DropChange) Hash(h hash.Hash) {
	h.Write(d.Account())
}

func (r RateVote) Hash(h hash.Hash) {
	h.Write(util.UInt16ToBA(uint16(r)))
}

func (k Key) Hash(h hash.Hash) {
	h.Write(k.X())
	h.Write(k.Y())
}

func (s Signature) Hash(h hash.Hash) {
	h.Write(s.R())
	h.Write(s.S())
}

func (ks KeySignature) Hash(h hash.Hash) {
	ks.Key().Hash(h)
	ks.Signature().Hash(h)
}

func (a Authorization) Hash(h hash.Hash) {
	a.Signatures().At(0).Hash(h)
}
