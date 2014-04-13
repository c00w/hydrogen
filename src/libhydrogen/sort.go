package libhydrogen

import (
	"libhydrogen/message"
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
