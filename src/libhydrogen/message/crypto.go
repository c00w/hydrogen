package message

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"hash"
	"math/big"
)

func bigint(b []byte) *big.Int {
	return big.NewInt(0).SetBytes(b)
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
