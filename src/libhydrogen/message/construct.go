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
