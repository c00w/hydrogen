package libhydrogen

import (
	"errors"
	"log"
	"sort"
	"sync"
	"time"

	"libhydrogen/message"
	"util"
)

type Hydrogen struct {
	currentledger *Ledger

	votes      []message.Vote
	newvote    chan message.Vote
	votetiming map[string]time.Time

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
	return &Hydrogen{l, nil, make(chan message.Vote), make(map[string]time.Time),
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

func (h *Hydrogen) GetLedger() *Ledger {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return h.currentledger
}

func (h *Hydrogen) handleVote(v message.Vote) {
	h.lock.Lock()
	util.Debugf("vote recieved %v", v)
	h.votes = append(h.votes, v)
	h.votetiming[util.Hash(v)] = time.Now()
	h.lock.Unlock()
}

func (h *Hydrogen) handleChange(c message.Change) {
	h.lock.Lock()
	util.Debugf("change recieved %v", c)
	h.changes = append(h.changes, c)
	h.lock.Unlock()
}

func (h *Hydrogen) eventloop() {
	for {
		bt := <-h.blocktimer.Chan()
		h.lock.Lock()

		ratechange := h.calculateRateChange()

		appliedchanges, appliedvotes := h.applyVotes(bt)

		h.blocktimer.SetTau(h.currentledger.Tau)

		h.cleanupChanges(appliedchanges)
		h.cleanupVotes(bt)

		h.changes = append(h.changes, ratechange)

		vote := message.NewSignedVote(h.changes, h.mp.node.Key)
		h.lock.Unlock()

		h.mp.SendVote(vote)

		if h.newblock != nil {
			h.newblock <- appliedvotes
		}
	}
}

func (h *Hydrogen) calculateRateChange() message.Change {
	times := make([]time.Time, 0, len(h.votetiming))

	for _, t := range h.votetiming {
		times = append(times, t)
	}

	h.votetiming = make(map[string]time.Time)

	sort.Sort(earliest(times))

	if len(times) == 0 {
		return message.NewSignedRateChange(message.RATEVOTE_CONSTANT, h.mp.node.Key)
	}

	median := times[len(times)/2]

	estimatedtau := median.Sub(h.currentledger.Created)

	if estimatedtau > h.currentledger.Tau/2 {
		return message.NewSignedRateChange(message.RATEVOTE_INCREASE, h.mp.node.Key)
	}

	if estimatedtau < h.currentledger.Tau/4 {
		return message.NewSignedRateChange(message.RATEVOTE_DECREASE, h.mp.node.Key)
	}

	return message.NewSignedRateChange(message.RATEVOTE_CONSTANT, h.mp.node.Key)
}

func (h *Hydrogen) applyVotes(t TimeRange) ([]message.Change, []message.Vote) {

	changes := make(map[string]message.Change)
	changecount := make(map[string]uint)

	appliedvotes := make([]message.Vote, 0)
	for _, v := range h.votes {
		if v.Time().Time().After(t.Start) && v.Time().Time().Before(t.End) {
			appliedvotes = append(appliedvotes, v)
		}
	}

	for _, v := range appliedvotes {
		for _, c := range v.Votes().ToArray() {
			id := util.Hash(c)
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

	ledger := h.currentledger.Copy(t.End)

	sort.Sort(timesort(appliedchanges))
	for _, change := range appliedchanges {
		err := ledger.Apply(change)
		if err != nil {
			panic(err)
		}
	}

	h.currentledger = ledger

	return appliedchanges, appliedvotes
}

func (h *Hydrogen) cleanupChanges(applied []message.Change) {
	seen := make(map[string]bool)

	for _, v := range applied {
		seen[util.Hash(v)] = true
	}

	notseen := make([]message.Change, 0)

	for _, v := range h.changes {

		ok := seen[util.Hash(v)]
		if !ok {
			notseen = append(notseen, v)
		}
	}

	h.changes = notseen

	sort.Sort(timesort(h.changes))

	changes := make([]message.Change, 0)
	changeledger := h.currentledger.Copy(time.Now())

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
