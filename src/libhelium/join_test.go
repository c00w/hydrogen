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
	key3 := util.GenKey()

	n1 := libnode.NewNode(key1, "l1")
	n2 := libnode.NewNode(key2, "l2")
	n3 := libnode.NewNode(key3, "l3")

	l := libhydrogen.NewLedger()
	l.AddEntry(util.KeyString(key1), "l1", 100)
	l.AddEntry(util.KeyString(key2), "l2", 100)
	l.AddEntry(util.KeyString(key3), "", 100)

	h1 := libhydrogen.NewHydrogen(n1)
	h2 := libhydrogen.NewHydrogen(n2)
	h3 := libhydrogen.NewHydrogen(n3)

	h1.AddLedger(l)
	h2.AddLedger(l)

	NewServer(n1, h1)

	n1.Listen("localhost:4010")
	n2.Connect("localhost:4010", "hydrogen")

	h1.WaitNewLedger()

	nl := Connect(n3, "localhost:4010")

	h3.AddLedger(nl)
	h3.WaitNewLedger()

}
