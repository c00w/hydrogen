package libhydrogen

import (
	"crypto/ecdsa"
	"crypto/sha512"

	"libhydrogen/message"
)

type Ledger struct {
	Accounts map[string]*Account
}

func NewLedger() *Ledger {
	return &Ledger{make(map[string]*Account)}
}

func (l *Ledger) Verify(auth message.Authorization, hash []byte) bool {

	if _, ok := l.Accounts[auth.Account()]; !ok {
		return false
	}

	ks := auth.Signatures().At(0)

	h := sha512.New()
	ks.Key().Hash(h)
	if string(h.Sum(nil)) != l.Accounts[auth.Account()].Key {
		return false
	}

	key := ks.Key().ECDSA()
	r, s := ks.Signature().Parse()
	return ecdsa.Verify(key, hash, r, s)
}

func (l *Ledger) AddEntry(account string, key string, location string) {
	l.Accounts[account] = &Account{account, key, location, 0}
}
