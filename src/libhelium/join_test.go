package libhelium

import (
	"testing"

	"libhydrogen"
	"libnode"
	"util"
)

func TestJoin(t *testing.T) {
	key1 := util.GenKey()
	key2 := util.GenKey()

	n1 := libnode.NewNode(key1, "l1")
	n2 := libnode.NewNode(key2, "l2")

	l := libhydrogen.NewLedger()
	l.AddEntry(util.KeyString(key1), "l1", 100)
	l.AddEntry(util.KeyString(key2), "", 100)

	h1 := libhydrogen.NewHydrogen(n1)
	h2 := libhydrogen.NewHydrogen(n2)

	h1.AddLedger(l)

	NewServer(n1, h1)

	n1.Listen("localhost:4010")
	n2.Connect("localhost:4010", "hydrogen")

	h1.WaitNewLedger()
	h1.WaitNewLedger()

	nl := Connect(n2, "localhost:4010")

	h2.AddLedger(nl)
	h2.WaitNewLedger()
	h2.WaitNewLedger()

}
