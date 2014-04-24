package libhydrogen

import (
	"sync"
	"time"
)

type TimeRange struct {
	Start time.Time
	End   time.Time
}

func NewBlockTimer(tau time.Duration, lasttime time.Time) *BlockTimer {
	b := &BlockTimer{tau, lasttime, nil, make(chan TimeRange), &sync.RWMutex{}}
	b.setupTimer()
	return b
}

type BlockTimer struct {
	tau          time.Duration
	lasttime     time.Time
	currenttimer *time.Timer
	firingchan   chan TimeRange
	lock         *sync.RWMutex
}

func (b *BlockTimer) Chan() chan TimeRange {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.firingchan
}

func (b *BlockTimer) SetTau(tau time.Duration) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if tau == b.tau {
		return
	}
	b.tau = tau
	b.currenttimer.Stop()
	b.setupTimer()
}

func (b *BlockTimer) setupTimer() {
	if b.lasttime.Add(b.tau).Before(time.Now()) {
		go b.timerFired()
		return
	}
	wait := time.Duration(b.tau.Nanoseconds() - (time.Now().Sub(b.lasttime).Nanoseconds()))
	b.currenttimer = time.AfterFunc(wait, b.timerFired)
}

func (b *BlockTimer) timerFired() {
	b.lock.Lock()
	b.firingchan <- TimeRange{b.lasttime, b.lasttime.Add(b.tau)}
	b.lasttime = b.lasttime.Add(b.tau)
	b.setupTimer()
	b.lock.Unlock()
}
