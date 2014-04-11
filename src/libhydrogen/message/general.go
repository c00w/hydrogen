package message

import (
	"hash"
)

type Verifier interface {
	Verify(ks Authorization, hash []byte) error
}

func (m Message) Verify(l Verifier, h hash.Hash) error {
	m.Payload().Hash(h)
	for _, ks := range m.AuthChain().ToArray() {
        if err := l.Verify(ks, h.Sum(nil)); err != nil {
			return err
		}
		ks.Hash(h)
	}

	return nil
}

func (m MessagePayload) Hash(h hash.Hash) {
	switch m.Which() {
	case MESSAGEPAYLOAD_VOTE:
		m.Vote().Hash(h)
	case MESSAGEPAYLOAD_CHANGE:
		m.Change().Hash(h)
	default:
	}
}

func (t Time) Hash(h hash.Hash) {
	h.Write(uint64toba(t.Seconds()))
	h.Write(uint32toba(t.NanoSeconds()))
}

func (v Vote) Hash(h hash.Hash) {
	for _, c := range v.Votes().ToArray() {
		c.Hash(h)
	}

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
	case CHANGETYPE_KEY:
		c.Type().Key().Hash(h)
	case CHANGETYPE_DROP:
		c.Type().Drop().Hash(h)
	case CHANGETYPE_TIME:
		c.Type().Time().Hash(h)
	default:
	}
}

func (t TransactionChange) Hash(h hash.Hash) {
	h.Write(t.Source())
	h.Write(t.Destination())
	h.Write(uint64toba(t.Amount()))
}

func (l LocationChange) Hash(h hash.Hash) {
	h.Write(l.Account())
	h.Write([]byte(l.Location()))
}

func (k KeyChange) Hash(h hash.Hash) {
	h.Write(k.Account())
	h.Write(k.Newkeys().At(0))
}

func (d DropChange) Hash(h hash.Hash) {
	h.Write(d.Account())
}

func (r RateChange) Hash(h hash.Hash) {
	h.Write(uint16toba(uint16(r.Vote())))
}
