package message

import (
	"crypto/ecdsa"
	"crypto/sha512"
	"hash"
	"time"

	"util"

	capnp "github.com/glycerine/go-capnproto"
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

func NewTimeNow(n *capnp.Segment) Time {
	t := NewTime(n)
	t.SetTime(time.Now())
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

func NewSignedTransaction(key *ecdsa.PrivateKey, destination string, amount uint64) Change {
	n := capnp.NewBuffer(nil)

	c := NewRootChange(n)
	c.SetCreated(NewTimeNow(n))

	t := NewTransactionChange(n)
	t.SetSource([]byte(util.KeyString(key)))
	t.SetDestination([]byte(destination))
	t.SetAmount(amount)

	c.Type().SetTransaction(t)

	s := sha512.New()
	c.Created().Hash(s)
	t.Hash(s)

	auth := NewSignedAuthorization(n, key, s.Sum(nil))
	c.SetAuthorization(auth)
	return c
}
