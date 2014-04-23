package libhydrogen

import (
	"testing"

	"libhydrogen/message"
	"libnode"
	"util"
)

func TestHydrogen(t *testing.T) {
	key1 := util.GenKey()
	key2 := util.GenKey()

	n1 := libnode.NewNode(key1, "location1")
	n2 := libnode.NewNode(key2, "location2")

	l := NewLedger()
	l.AddEntry(util.KeyString(key1), "location1", 100)
	l.AddEntry(util.KeyString(key2), "location2", 100)

	tc1 := make(chan []message.Vote)
	tc2 := make(chan []message.Vote)

	h1 := newHydrogen(n1, tc1)
	h2 := newHydrogen(n2, tc2)

	h1.AddLedger(l)
	h2.AddLedger(l)

	n1.Listen("localhost:4005")
	n2.Connect("localhost:4005", "hydrogen")

	node2 := util.KeyString(key2)

	h1.TransferMoney(node2, 10)

	<-tc1
	<-tc2

	v1 := <-tc1
	v2 := <-tc2

	if len(v1) != 2 {
		for _, v := range v1 {
			t.Logf("%v", v)
		}
		t.Errorf("Not enough votes %d", len(v1))
	}

	if len(v2) != 2 {
		for _, v := range v2 {
			t.Logf("%v", v)
		}
		t.Errorf("Not enough votes %d", len(v2))
	}

	if b, err := h1.GetBalance(node2); b != 110 {
		if err != nil {
			t.Errorf("error fetching balance", err)
		}
		t.Logf("node1 balance is %d", b)
		t.Errorf("node2 balance is %d != 110", b)
	}

	if b, err := h2.GetBalance(node2); b != 110 {
		if err != nil {
			t.Errorf("error fetching balance", err)
		}
		t.Logf("node1 balance is %d", b)
		t.Errorf("node2 balance is %d != 110", b)
	}

}
