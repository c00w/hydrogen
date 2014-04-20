package message

import (
	"crypto/ecdsa"

	"util"

	capnp "github.com/glycerine/go-capnproto"
)

func NewSignedRateChange(r RateVote, k *ecdsa.PrivateKey) Change {

	n := capnp.NewBuffer(nil)

	c := NewRootChange(n)
	c.SetCreated(NewTimeNow(n))

	rc := NewRateChange(n)
	rc.SetVote(r)

	c.Type().SetTime(rc)

	h := util.Hash(c.Created(), c.Type().Time())

	c.SetAuthorization(NewSignedAuthorization(n, k, []byte(h)))
	return c
}

func NewSignedVote(c []Change, key *ecdsa.PrivateKey) Vote {
	ns := capnp.NewBuffer(nil)

	v := NewRootVote(ns)
	cl := NewChangeList(ns, len(c))
	for i, v := range c {
		capnp.PointerList(cl).Set(i, capnp.Object(v))
	}
	v.SetVotes(cl)
	v.SetTime(NewTimeNow(ns))

	h := util.Hash(v.Votes(), v.Time())

	a := NewSignedAuthorization(ns, key, []byte(h))
	v.SetAuthorization(a)
	return v
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

	h := util.Hash(c.Created(), c.Type().Transaction())
	auth := NewSignedAuthorization(n, key, []byte(h))
	c.SetAuthorization(auth)
	return c
}
