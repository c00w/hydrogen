package libhydrogen

import (
	"crypto/sha512"
	"log"
	"sort"
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

    newblock chan struct{}
}

func NewHydrogen(l *Ledger, b *BlockTimer) *Hydrogen {
    h := &Hydrogen{l, nil, make(chan message.Vote),
                        nil, make(chan message.Change), b, nil, nil}
    go h.eventloop()
    return h
}

func (h *Hydrogen) RegisterBus(mp *MessagePasser) {
    h.mp = mp
}

func (h *Hydrogen) Verify(ks message.Authorization, hash []byte) error {
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
	h.newvote <- v

}

func (h *Hydrogen) handleChange(c message.Change) {
	h.newchange <- c
}

func (h *Hydrogen) eventloop() {
	for {
		select {
		case c := <-h.newchange:
			h.changes = append(h.changes, c)

		case v := <-h.newvote:
			h.votes = append(h.votes, v)

		case <-h.blocktimer.Chan():
			newledger, appliedchanges, _ := h.tallyVotes()
			h.currentledger = newledger
			h.changes = h.filterChanges(appliedchanges)
			h.changes = h.validateChanges()
			vote := h.createVote()
			h.mp.SendVote(vote)
            select {
            case h.newblock <- struct{}{}:
            default:
            }
		}
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

	cl := message.NewChangeList(ns, len(h.changes))
	for i, v := range h.changes {
		capnp.PointerList(cl).Set(i, capnp.Object(v))
	}

	v := message.NewVote(ns)
	v.SetVotes(cl)
	v.Time().SetTime(time.Now())

	v.Authorization()

	return v
}
