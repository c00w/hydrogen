package message

import (
	"crypto/ecdsa"
	"hash"

	"util"

	capnp "github.com/glycerine/go-capnproto"
)

func NewSignedVote(c []Change, r RateVote, key *ecdsa.PrivateKey) Vote {
	ns := capnp.NewBuffer(nil)

	v := NewRootVote(ns)
	cl := NewChangeList(ns, len(c))
	for i, v := range c {
		capnp.PointerList(cl).Set(i, capnp.Object(v))
	}
	v.SetVotes(cl)
	v.SetTime(util.NewTimeNow(ns))
	v.SetRate(r)

	h := util.Hash(v.Votes(), v.Rate(), v.Time())

	a := NewSignedAuthorization(ns, key, []byte(h))
	v.SetAuthorization(a)
	return v
}

func NewSignedTransaction(key *ecdsa.PrivateKey, destination string, amount uint64) Change {
	n := capnp.NewBuffer(nil)

	c := NewRootChange(n)
	c.SetCreated(util.NewTimeNow(n))

	t := NewTransactionChange(n)
	t.SetSource([]byte(util.KeyString(key)))
	t.SetDestination([]byte(destination))
	t.SetAmount(amount)

	c.Type().SetTransaction(t)

	h := util.Hash(c.Created(), c.Type().Transaction())
	auth := NewSignedAuthorization(n, key, []byte(h))
	c.SetAuthorization(auth)
	return c
}

func CreateMessageFromChange(c Change, key *ecdsa.PrivateKey) *capnp.Segment {
	n := capnp.NewBuffer(nil)

	m := NewRootMessage(n)
	m.Payload().SetChange(c)

	a := NewSignedAuthorization(n, key, []byte(util.Hash(c)))

	al := NewAuthorizationList(n, 1)
	capnp.PointerList(al).Set(0, capnp.Object(a))

	m.SetAuthChain(al)
	return n
}

func CreateMessageFromVote(v Vote, key *ecdsa.PrivateKey) *capnp.Segment {
	n := capnp.NewBuffer(nil)

	m := NewRootMessage(n)
	m.Payload().SetVote(v)

	a := NewSignedAuthorization(n, key, []byte(util.Hash(v)))

	al := NewAuthorizationList(n, 1)
	capnp.PointerList(al).Set(0, capnp.Object(a))

	m.SetAuthChain(al)
	return n
}

func AppendAuthMessage(m Message, run hash.Hash, key *ecdsa.PrivateKey) *capnp.Segment {

	n := capnp.NewBuffer(nil)

	m2 := NewRootMessage(n)

	l := NewAuthorizationList(n, m.AuthChain().Len()+1)
	for i, v := range m.AuthChain().ToArray() {
		capnp.PointerList(l).Set(i, capnp.Object(v))
	}

	a := NewSignedAuthorization(n, key, run.Sum(nil))

	capnp.PointerList(l).Set(m.AuthChain().Len(), capnp.Object(a))

	m2.SetAuthChain(l)

	switch m.Payload().Which() {
	case MESSAGEPAYLOAD_VOTE:
		m2.Payload().SetVote(m.Payload().Vote())
	case MESSAGEPAYLOAD_CHANGE:
		m2.Payload().SetChange(m.Payload().Change())
	}

	return n
}
