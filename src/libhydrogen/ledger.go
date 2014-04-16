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

func (l *Ledger) Apply(c message.Change) error {

	s := sha512.New()
	c.Created().Hash(s)

	switch c.Type().Which() {
	case message.CHANGETYPE_TRANSACTION:
		t := c.Type().Transaction()
		t.Hash(s)

		err := l.Verify(c.Authorization(), s.Sum(nil))
		if err != nil {
			return err
		}

		source, ok := l.Accounts[string(t.Source())]
		if !ok {
			return errors.New("no such source account")
		}

		destination, ok := l.Accounts[string(t.Destination())]
		if !ok {
			return errors.New("no such destination account")
		}

		if source.Balance < t.Amount() {
			return errors.New("insufficient funds")
		}

		source.Balance -= t.Amount()
		destination.Balance += t.Amount()

		l.Accounts[source.ID] = source
		l.Accounts[destination.ID] = destination

	case message.CHANGETYPE_LOCATION:
		lo := c.Type().Location()
		lo.Hash(s)

		err := l.Verify(c.Authorization(), s.Sum(nil))
		if err != nil {
			return err
		}

		account, ok := l.Accounts[string(lo.Account())]
		if !ok {
			return errors.New("no such account")
		}

		account.Location = lo.Location()
		l.Accounts[account.ID] = account

	case message.CHANGETYPE_KEY:
		k := c.Type().Key()
		k.Hash(s)

		err := l.Verify(c.Authorization(), s.Sum(nil))
		if err != nil {
			return err
		}

		account, ok := l.Accounts[string(k.Account())]
		if !ok {
			return errors.New("no such account")
		}

		account.Key = string(k.Newkeys().At(0))
		l.Accounts[account.ID] = account

	case message.CHANGETYPE_DROP:
		d := c.Type().Key()
		account := string(d.Account())
		info, ok := l.Accounts[account]
		if !ok {
			return errors.New("no such account")
		}
		info.Location = ""
		l.Accounts[account] = info

	case message.CHANGETYPE_TIME:

	default:
		return errors.New("unrecognized change type")
	}

	return nil
}

func (l *Ledger) HostCount() uint {
	i := uint(0)
	for _, a := range l.Accounts {
		if a.Location != "" {
			i += 1
		}
	}
	return i
}

func (l *Ledger) Copy() *Ledger {
	nl := &Ledger{make(map[string]*Account)}
	for k, v := range l.Accounts {
		nl.Accounts[k] = v
	}
	return nl
}
