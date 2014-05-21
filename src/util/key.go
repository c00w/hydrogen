package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
	"os"

	capnp "github.com/glycerine/go-capnproto"
)

// Create an encoded P521 key from a normal ecdsa private key
func NewEncodedP521Key(k *ecdsa.PrivateKey) *capnp.Segment {
	b := capnp.NewBuffer(nil)

	p := NewRootP521Key(b)
	p.SetD(k.D.Bytes())
	p.SetX(k.PublicKey.X.Bytes())
	p.SetY(k.PublicKey.Y.Bytes())

	return b
}

// Parse an encoded P521 key into a normal ecdsa private key
func (p P521Key) ParseKey() *ecdsa.PrivateKey {
	pub := ecdsa.PublicKey{
		elliptic.P521(),
		big.NewInt(0).SetBytes(p.X()),
		big.NewInt(0).SetBytes(p.Y()),
	}

	priv := &ecdsa.PrivateKey{pub, big.NewInt(0).SetBytes(p.D())}
	return priv
}

// Generate an encoded key at the given file location
func GenerateKey(path string) error {
	fd, err := os.Create(path)
	if err != nil {
		return err
	}

	key := GenKey()
	enc := NewEncodedP521Key(key)
	_, err = enc.WriteTo(fd)
	fd.Close()
	return err
}

// Load an encoded key from the given file location
func LoadKey(path string) (*ecdsa.PrivateKey, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	seg, err := capnp.ReadFromStream(fd, nil)
	if err != nil {
		return nil, err
	}

	key := ReadRootP521Key(seg)
	return key.ParseKey(), nil
}
