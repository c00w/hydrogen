package libhydrogen

import (
	"crypto/sha512"
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
		<-h.blocktimer.Chan()
        h.lock.Lock()
        oldvotes := h.votes
        newledger, appliedchanges, _ := h.tallyVotes()
        h.currentledger = newledger
        h.changes = h.filterChanges(appliedchanges)
        h.changes = h.validateChanges()
        vote := h.createVote()
        h.mp.SendVote(vote)
        if h.newblock != nil {
            h.newblock <- oldvotes
        }
        h.lock.Unlock()
	}
}

func (h *Hydrogen) tallyVotes() (*Ledger, []message.Change, time.Duration) {
	changes := make(map[string]message.Change)
	changecount := make(map[string]uint)

	s := sha512.New()

	for _, v := range h.votes {
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
		if count > h.currentledger.HostCount() {
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

	return ledger, appliedchanges, 0
}

func (h *Hydrogen) filterChanges(applied []message.Change) []message.Change {
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

	return notseen
}

func (h *Hydrogen) validateChanges() []message.Change {
	sort.Sort(timesort(h.changes))

	changes := make([]message.Change, 0)
	changeledger := h.currentledger.Copy()

	for _, change := range h.changes {
		err := changeledger.Apply(change)
		if err == nil {
			changes = append(changes, change)
		}
	}

	return changes
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

    v.SetAccount(h.mp.node.Account)

    s := sha512.New()
    cl.Hash(s)
    t.Hash(s)
    s.Write([]byte(h.mp.node.Account))

    a := message.NewSignedAuthorization(ns, h.mp.node.Account, h.mp.node.Key, s.Sum(nil))
	v.SetAuthorization(a)
	return v
}
