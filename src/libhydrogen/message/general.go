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
