package libhydrogen

import (
	"testing"
	"time"
)

func TestBlockTimer(t *testing.T) {
	now := time.Now()
	tau := 100 * time.Millisecond

	b := NewBlockTimer(100*time.Millisecond, now)

	<-b.Chan()

	now2 := time.Now()
	if now2.Before(now.Add(tau)) {
		t.Fatalf("blocktimer returned to quickly %s < %s", now2.Sub(now), tau)
	}
}
