package libhydrogen

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"time"

	"libhydrogen/message"
	"util"
)

var ZEROACCOUNT string = util.Hash()

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

func (l *Ledger) Apply(c message.Change, votesseen map[string]bool) error {

	switch c.Type().Which() {
	case message.CHANGETYPE_TRANSACTION:
		t := c.Type().Transaction()
		h := util.Hash(c.Created(), t)

		err := l.Verify(c.Authorization(), []byte(h))
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
		h := util.Hash(c.Created(), lo)

		err := l.Verify(c.Authorization(), []byte(h))
		if err != nil {
			return err
		}

		account, ok := l.Accounts[string(lo.Account())]
		if !ok {
			return errors.New("no such account")
		}

		if votesseen != nil && !votesseen[string(lo.Account())] {
			return errors.New("account is not active")
		}

		account.Location = lo.Location()
		l.Accounts[account.Key] = account

	default:
		return errors.New("unrecognized change type")
	}

	return nil
}

func (l *Ledger) ApplyRate(r message.RateVote) {
	switch r {
	case message.RATEVOTE_INCREASE:
		l.Tau = l.Tau * 11 / 10
	case message.RATEVOTE_DECREASE:
		l.Tau = l.Tau * 10 / 11
	default:
	}
}

func calculateLoss(amount, period uint64) uint64 {
	return (amount * period >> 47) + 1
}

func (l *Ledger) ApplyWealthRedistribution(period time.Duration) {

	totaltaken := uint64(0)
	hostcount := l.HostCount()
	if hostcount == 0 {
		hostcount = 1
	}

	expected := uint64(1) << 63 / hostcount

	for k, account := range l.Accounts {
		allowed := uint64(0)
		if account.Active() {
			allowed = expected
		}
		if account.Balance > allowed {
			loss := calculateLoss(account.Balance-allowed, uint64(period))
			totaltaken += loss
			account := account.Copy()
			account.Balance -= loss
			l.Accounts[k] = account
		}
	}

	redistribution := totaltaken / hostcount

	for k, account := range l.Accounts {
		if account.Active() {
			account = account.Copy()
			account.Balance += redistribution
			l.Accounts[k] = account
		}
	}

	extra := totaltaken - redistribution*hostcount

	c, ok := l.Accounts[ZEROACCOUNT]
	if !ok {
		c = &Account{ZEROACCOUNT, "", 0}
	}

	c = c.Copy()
	c.Balance += extra
	l.Accounts[ZEROACCOUNT] = c

}

func (l *Ledger) Drop(n string) {
	a, ok := l.Accounts[n]
	if !ok {
		return
	}
	a = a.Copy()
	a.Location = ""
	l.Accounts[n] = a
}

func (l *Ledger) HostCount() uint64 {
	i := uint64(0)
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
