package libhydrogen

import (
	"libhydrogen/message"
	"time"
)

type timesort []message.Change

func (t timesort) Len() int {
	return len(t)
}

func (t timesort) Less(i, j int) bool {
	return t[i].Created().Time().Before(t[j].Created().Time())
}

func (t timesort) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

type earliest []time.Time

func (t earliest) Len() int {
	return len(t)
}

func (t earliest) Less(i, j int) bool {
	return t[i].Before(t[j])
}

func (t earliest) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
