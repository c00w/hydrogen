package libhydrogen

import (
	"crypto/ecdsa"
	"crypto/sha512"
	"errors"
	"fmt"
	"time"

	"libhydrogen/message"
)

type Ledger struct {
	Accounts map[string]*Account
	Tau      time.Duration
	Created  time.Time
}

func NewLedger() *Ledger {
	return &Ledger{make(map[string]*Account), time.Second, time.Now()}
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

func (l *Ledger) AddEntry(key string, location string, balance uint64) {
	l.Accounts[key] = &Account{key, location, balance}
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
		source = source.Copy()

		destination, ok := l.Accounts[string(t.Destination())]
		if !ok {
			l.AddEntry(string(t.Destination()), "", 0)
			destination = l.Accounts[string(t.Destination())]
		}
		destination = destination.Copy()

		if source.Balance < t.Amount() {
			return errors.New("insufficient funds")
		}

		source.Balance -= t.Amount()
		destination.Balance += t.Amount()

		l.Accounts[source.Key] = source
		l.Accounts[destination.Key] = destination

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
		l.Accounts[account.Key] = account

	case message.CHANGETYPE_DROP:
		d := c.Type().Drop()
		account := string(d.Account())
		info, ok := l.Accounts[account]
		if !ok {
			return errors.New("no such account")
		}
		info.Location = ""
		l.Accounts[account] = info

	case message.CHANGETYPE_TIME:
		switch c.Type().Time().Vote() {
		case message.RATEVOTE_INCREASE:
			l.Tau = l.Tau * 11 / 10
		case message.RATEVOTE_DECREASE:
			l.Tau = l.Tau * 10 / 11
		default:
		}

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

func (l *Ledger) Copy(t time.Time) *Ledger {
	nl := &Ledger{make(map[string]*Account), l.Tau, t}
	for k, v := range l.Accounts {
		nl.Accounts[k] = v
	}
	return nl
}
