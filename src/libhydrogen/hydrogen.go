package libhydrogen

import (
	"errors"
	"log"
	"sort"
	"sync"
	"time"

	"libhydrogen/message"
	"libnode"
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

	mp *messagePasser

	newblock chan []message.Vote

	lock *sync.RWMutex

	newledger *sync.Cond

	disabledrop bool
}

func NewHydrogen(n *libnode.Node) *Hydrogen {
	return newHydrogen(n, nil)
}

func newHydrogen(n *libnode.Node, c chan []message.Vote) *Hydrogen {
	h := &Hydrogen{}
	h.newvote = make(chan message.Vote)
	h.votetiming = make(map[string]time.Time)
	h.newchange = make(chan message.Change)
	h.mp = newMessagePasser(n, h)
	h.newblock = c
	h.lock = &sync.RWMutex{}
	h.newledger = sync.NewCond(h.lock)
	return h
}

func (h *Hydrogen) Verify(ks message.Authorization, hash []byte) error {
	h.lock.RLock()
	defer h.lock.RUnlock()
	if h.currentledger == nil {
		return nil
	}
	return h.currentledger.Verify(ks, hash)
}

func (h *Hydrogen) shost() string {
	return util.Short(util.KeyString(h.mp.node.Key))
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

func (h *Hydrogen) SetLocation(newlocation string) error {
	u := message.NewSignedLocationUpdate(h.mp.node.Key, newlocation)
	h.mp.SendChange(u)
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

func (h *Hydrogen) GetLocation(account string) (string, error) {
	h.lock.RLock()
	entry, ok := h.currentledger.Accounts[account]
	h.lock.RUnlock()
	if !ok {
		return "", errors.New("no such account")
	}
	return entry.Location, nil
}

func (h *Hydrogen) WaitNewLedger() *Ledger {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.newledger.Wait()
	return h.currentledger
}

func (h *Hydrogen) AddLedger(l *Ledger) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.currentledger != nil {
		panic("Add ledger called twice...")
	}
	h.currentledger = l
	h.blocktimer = NewBlockTimer(l.Tau, l.Created)

	t := TimeRange{l.Created.Add(-l.Tau), l.Created}

	h.cleanupVotes(t)

	go h.eventloop()
}

func (h *Hydrogen) GetLedger() *Ledger {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return h.currentledger
}

func (h *Hydrogen) handleVote(v message.Vote) {
	h.lock.Lock()
	util.Debugf("Host %s, vote recieved %v", h.shost(), v)
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
	util.Debugf("Host %s starting eventloop", h.shost())
	for {
		bt := <-h.blocktimer.Chan()
		h.lock.Lock()

		appliedchanges, appliedvotes := h.applyVotes(bt)

		h.blocktimer.SetTau(h.currentledger.Tau)

		h.cleanupChanges(appliedchanges)

		if h.caughtup(bt) {
			util.Debugf("Host %s Generating vote", h.shost())
			vote := message.NewSignedVote(
				h.changes,
				h.calculateRateChange(),
				h.calculateDropVotes(),
				h.mp.node.Key)
			go h.mp.SendVote(vote)
		} else {
			util.Debugf("Not caughtup fastforwarding")
		}

		h.cleanupVotes(bt)

		h.lock.Unlock()

		if h.newblock != nil {
			h.newblock <- appliedvotes
		}
	}
}

func (h *Hydrogen) calculateRateChange() message.RateVote {
	times := make([]time.Time, 0, len(h.votetiming))

	for _, t := range h.votetiming {
		times = append(times, t)
	}

	h.votetiming = make(map[string]time.Time)

	sort.Sort(earliest(times))

	if len(times) == 0 || len(times) == 1 {
		return message.RATEVOTE_CONSTANT
	}

	median := times[len(times)/2]

	estimatedtau := median.Sub(h.currentledger.Created)

	if estimatedtau > h.currentledger.Tau/2 {
		return message.RATEVOTE_INCREASE
	}

	if estimatedtau < h.currentledger.Tau/4 {
		return message.RATEVOTE_DECREASE
	}

	return message.RATEVOTE_CONSTANT
}

func (h *Hydrogen) calculateDropVotes() message.DropChange_List {
	if h.disabledrop {
		h.disabledrop = false
		return message.NewDropVotes(nil)
	}

	absenthosts := make([]string, 0)

	seen := make(map[string]bool)

	for _, v := range h.votes {
		seen[v.Authorization().Account()] = true
	}

	for host, a := range h.currentledger.Accounts {
		if a.Active() && !seen[host] {
			absenthosts = append(absenthosts, host)
		}
	}

	return message.NewDropVotes(absenthosts)
}

func (h *Hydrogen) applyVotes(t TimeRange) ([]message.Change, []message.Vote) {

	changes := make(map[string]message.Change)
	changecount := make(map[string]uint64)
	drop := make(map[string]uint64)

	appliedvotes := make([]message.Vote, 0)
	votesseen := make(map[string]bool)
	for _, v := range h.votes {
		if v.Time().Time().After(t.Start) && v.Time().Time().Before(t.End) {
			appliedvotes = append(appliedvotes, v)
			votesseen[v.Authorization().Account()] = true
		}
	}

	faster := uint64(0)
	slower := uint64(0)

	for _, v := range appliedvotes {
		switch v.Rate() {
		case message.RATEVOTE_INCREASE:
			faster += 1
		case message.RATEVOTE_DECREASE:
			slower += 1
		}

		for _, d := range v.Drop().ToArray() {
			drop[string(d.Account())] += 1
		}
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

	if faster > h.currentledger.HostCount()/2 {
		ledger.ApplyRate(message.RATEVOTE_INCREASE)
	}

	if slower > h.currentledger.HostCount()/2 {
		ledger.ApplyRate(message.RATEVOTE_DECREASE)
	}

	sort.Sort(timesort(appliedchanges))
	for _, change := range appliedchanges {
		err := ledger.Apply(change, votesseen)
		if err != nil {
			panic(err)
		}
	}

	for account, c := range drop {
		if c > h.currentledger.HostCount()/2 {
			ledger.Drop(account)
		}
	}

	ledger.ApplyWealthRedistribution(t.End.Sub(t.Start))

	h.currentledger = ledger
	h.newledger.Broadcast()

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
		err := changeledger.Apply(change, nil)
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

func (h *Hydrogen) caughtup(t TimeRange) bool {
	return t.End.Add(h.currentledger.Tau).After(time.Now())
}
