package libhydrogen

import (
	"crypto/ecdsa"
	"crypto/sha512"
    "errors"
    "fmt"

	"libhydrogen/message"
)

type Ledger struct {
	Accounts map[string]*Account
}

func NewLedger() *Ledger {
	return &Ledger{make(map[string]*Account)}
}

func (l *Ledger) Verify(auth message.Authorization, hash []byte) error {

	if _, ok := l.Accounts[auth.Account()]; !ok {
		return errors.New(fmt.Sprintf("no such account in ledger acct=\"%s\"", auth.Account()))
	}

	ks := auth.Signatures().At(0)

	h := sha512.New()
	ks.Key().Hash(h)
	if string(h.Sum(nil)) != l.Accounts[auth.Account()].Key {
		return errors.New("invalid key")
	}

	key := ks.Key().ECDSA()
	r, s := ks.Signature().Parse()
    ok := ecdsa.Verify(key, hash, r, s)
    if !ok {
        return errors.New("ecdsa verification failed")
    }
    return nil
}

func (l *Ledger) AddEntry(account string, key string, location string) {
	l.Accounts[account] = &Account{account, key, location, 0}
}
