package message

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"hash"
	"math/big"

	capnp "github.com/glycerine/go-capnproto"
)

func bigint(b []byte) *big.Int {
	return big.NewInt(0).SetBytes(b)
}

func NewKeyFromECDSA(n *capnp.Segment, o *ecdsa.PublicKey) Key {
	k := NewKey(n)
	k.SetX(o.X.Bytes())
	k.SetY(o.Y.Bytes())
	return k
}

func (k Key) ECDSA() *ecdsa.PublicKey {
	return &ecdsa.PublicKey{elliptic.P521(), bigint(k.X()), bigint(k.Y())}
}

func (k Key) Hash(h hash.Hash) {
	h.Write(k.X())
	h.Write(k.Y())
}

func (s Signature) Parse() (*big.Int, *big.Int) {
	return bigint(s.R()), bigint(s.S())
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
	h.Write([]byte(a.Account()))
	a.Signatures().At(0).Hash(h)
}

func NewSignedAuthorization(n *capnp.Segment, account string, key *ecdsa.PrivateKey, item []byte) Authorization {

	k := NewKey(n)
	k.SetX(key.X.Bytes())
	k.SetY(key.Y.Bytes())

	r, s, err := ecdsa.Sign(rand.Reader, key, item)
	if err != nil {
		panic(err)
	}

	sig := NewSignature(n)
	sig.SetR(r.Bytes())
	sig.SetS(s.Bytes())

	keysig := NewKeySignature(n)
	keysig.SetKey(k)
	keysig.SetSignature(sig)

	sl := NewKeySignatureList(n, 1)
	capnp.PointerList(sl).Set(0, capnp.Object(keysig))

	auth := NewAuthorization(n)
	auth.SetAccount(account)
	auth.SetSignatures(sl)

	return auth
}
