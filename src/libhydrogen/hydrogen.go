package libhydrogen

import (
	"crypto/sha512"
	"errors"
	"log"
	"sort"
	"sync"
	"time"

	"libhydrogen/message"

	capnp "github.com/glycerine/go-capnproto"
)

type Hydrogen struct {
	currentledger *Ledger

	votes   []message.Vote
	newvote chan message.Vote

	changes   []message.Change
	newchange chan message.Change

	blocktimer *BlockTimer

	mp *MessagePasser

	newblock chan []message.Vote

	lock *sync.RWMutex
}

func NewHydrogen(l *Ledger, b *BlockTimer) *Hydrogen {
	return newHydrogen(l, b, nil)
}

func newHydrogen(l *Ledger, b *BlockTimer, c chan []message.Vote) *Hydrogen {
	return &Hydrogen{l, nil, make(chan message.Vote),
		nil, make(chan message.Change), b, nil, c, &sync.RWMutex{}}
}

func (h *Hydrogen) RegisterBus(mp *MessagePasser) {
	h.mp = mp
	go h.eventloop()
}

func (h *Hydrogen) Verify(ks message.Authorization, hash []byte) error {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return h.currentledger.Verify(ks, hash)
}

func (h *Hydrogen) Handle(m message.Message) {
	switch m.Payload().Which() {
	case message.MESSAGEPAYLOAD_VOTE:
		h.handleVote(m.Payload().Vote())
	case message.MESSAGEPAYLOAD_CHANGE:
		h.handleChange(m.Payload().Change())
	default:
		log.Print("unknown message payload type")
	}
}

func (h *Hydrogen) TransferMoney(destination string, amount uint64) error {
	t := message.NewSignedTransaction(h.mp.node.Key, destination, amount)
	h.mp.SendChange(t)
	return nil
}

func (h *Hydrogen) GetBalance(account string) (uint64, error) {
	h.lock.RLock()
	entry, ok := h.currentledger.Accounts[account]
	h.lock.RUnlock()
	if !ok {
		return 0, errors.New("no such account")
	}
	return entry.Balance, nil
}

func (h *Hydrogen) handleVote(v message.Vote) {
	h.lock.Lock()
	h.votes = append(h.votes, v)
	h.lock.Unlock()
}

func (h *Hydrogen) handleChange(c message.Change) {
	h.lock.Lock()
	h.changes = append(h.changes, c)
	h.lock.Unlock()
}

func (h *Hydrogen) eventloop() {
	for {
		bt := <-h.blocktimer.Chan()
		h.lock.Lock()

		appliedchanges, appliedvotes, _ := h.applyVotes(bt)

		h.cleanupChanges(appliedchanges)
		h.cleanupVotes(bt)

		vote := h.createVote()
		h.lock.Unlock()

		h.mp.SendVote(vote)

		if h.newblock != nil {
			h.newblock <- appliedvotes
		}
	}
}

func (h *Hydrogen) applyVotes(t TimeRange) ([]message.Change, []message.Vote, time.Duration) {

	changes := make(map[string]message.Change)
	changecount := make(map[string]uint)

	s := sha512.New()

	appliedvotes := make([]message.Vote, 0)
	for _, v := range h.votes {
		if v.Time().Time().After(t.Start) && v.Time().Time().Before(t.End) {
			appliedvotes = append(appliedvotes, v)
		}
	}

	for _, v := range appliedvotes {
		for _, c := range v.Votes().ToArray() {
			s.Reset()
			c.Hash(s)
			id := string(s.Sum(nil))
			if _, ok := changes[id]; !ok {
				changes[id] = c
			}
			changecount[id] = changecount[id] + 1
		}
	}

	appliedchanges := make([]message.Change, 0)

	for id, count := range changecount {
		if count > h.currentledger.HostCount()/2 {
			appliedchanges = append(appliedchanges, changes[id])
		}
	}

	ledger := h.currentledger.Copy()

	sort.Sort(timesort(appliedchanges))
	for _, change := range appliedchanges {
		err := ledger.Apply(change)
		if err != nil {
			panic(err)
		}
	}

	h.currentledger = ledger

	return appliedchanges, appliedvotes, 0
}

func (h *Hydrogen) cleanupChanges(applied []message.Change) {
	seen := make(map[string]bool)

	s := sha512.New()

	for _, v := range applied {
		s.Reset()
		v.Hash(s)
		seen[string(s.Sum(nil))] = true
	}

	notseen := make([]message.Change, 0)

	for _, v := range h.changes {
		s.Reset()
		v.Hash(s)

		ok := seen[string(s.Sum(nil))]
		if !ok {
			notseen = append(notseen, v)
		}
	}

	h.changes = notseen

	sort.Sort(timesort(h.changes))

	changes := make([]message.Change, 0)
	changeledger := h.currentledger.Copy()

	for _, change := range h.changes {
		err := changeledger.Apply(change)
		if err == nil {
			changes = append(changes, change)
		}
	}

	h.changes = changes
}

func (h *Hydrogen) cleanupVotes(t TimeRange) {
	votes := make([]message.Vote, 0)

	for _, v := range h.votes {
		if v.Time().Time().After(t.End) {
			votes = append(votes, v)
		}
	}

	h.votes = votes
}

func (h *Hydrogen) createVote() message.Vote {
	ns := capnp.NewBuffer(nil)

	v := message.NewRootVote(ns)
	cl := message.NewChangeList(ns, len(h.changes))
	for i, v := range h.changes {
		capnp.PointerList(cl).Set(i, capnp.Object(v))
	}
	v.SetVotes(cl)

	t := message.NewTime(ns)
	t.SetTime(time.Now())
	v.SetTime(t)

	s := sha512.New()
	cl.Hash(s)
	t.Hash(s)
	s.Write([]byte(h.mp.node.Account))

	a := message.NewSignedAuthorization(ns, h.mp.node.Key, s.Sum(nil))
	v.SetAuthorization(a)
	return v
}
