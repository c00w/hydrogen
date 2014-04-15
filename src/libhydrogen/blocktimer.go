package libhydrogen

import (
	"sync"
	"time"
)

func NewBlockTimer(tau time.Duration, lasttime time.Time) *BlockTimer {
	b := &BlockTimer{tau, lasttime, nil, make(chan struct{}), &sync.RWMutex{}}
	b.setupTimer()
	return b
}

type BlockTimer struct {
	tau          time.Duration
	lasttime     time.Time
	currenttimer *time.Timer
	firingchan   chan struct{}
	lock         *sync.RWMutex
}

func (b *BlockTimer) Chan() chan struct{} {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.firingchan
}

func (b *BlockTimer) IncreaseTau() {
	b.lock.Lock()
	b.tau = time.Duration(b.tau.Nanoseconds() * 11 / 10)
	b.currenttimer.Stop()
	b.setupTimer()
	b.lock.Unlock()
}

func (b *BlockTimer) DecreaseTau() {
	b.lock.Lock()
	b.tau = time.Duration(b.tau.Nanoseconds() * 10 / 11)
	b.currenttimer.Stop()
	b.setupTimer()
	b.lock.Unlock()
}

func (b *BlockTimer) setupTimer() {
	wait := time.Duration(b.tau.Nanoseconds() - (time.Now().Sub(b.lasttime).Nanoseconds() % b.tau.Nanoseconds()))
	b.currenttimer = time.AfterFunc(wait, b.timerFired)
}

func (b *BlockTimer) timerFired() {
	b.lock.Lock()
	b.lasttime = b.lasttime.Add(b.tau)
	b.firingchan <- struct{}{}
	b.setupTimer()
	b.lock.Unlock()
}
