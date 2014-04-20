// A utility package to make common things easier
package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha512"
	"hash"
)

func KeyString(k *ecdsa.PrivateKey) string {
	s := sha512.New()
	s.Write(k.X.Bytes())
	s.Write(k.Y.Bytes())
	return string(s.Sum(nil))
}

func GenKey() *ecdsa.PrivateKey {
	k, e := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if e != nil {
		panic(e)
	}
	return k
}

type Hashable interface {
	Hash(h hash.Hash)
}

func Hash(items ...Hashable) string {
	s := sha512.New()
	for _, item := range items {
		item.Hash(s)
	}
	return string(s.Sum(nil))
}
