package libberyllium

import (
	"testing"

	"libhydrogen"
	"libnode"
	"util"
)

func TestDownload(t *testing.T) {

	key1 := util.GenKey()
	key2 := util.GenKey()

	l := libhydrogen.NewLedger()
	l.AddEntry(util.KeyString(key1), "127.0.0.1:3012", 100)

	n1 := libnode.NewNode(key1, "127.0.0.1:3012")
	n2 := libnode.NewNode(key2, "127.0.0.1:3013")

	h1 := libhydrogen.NewHydrogen(n1)
	h1.SetBootStrap()
	h1.AddLedger(l)

	NewServer(n1, h1)

	n1.Listen("127.0.0.1:3012")

	err := GetMoney(n2, "127.0.0.1:3012", util.KeyString(key2))
	if err != nil {
		t.Fatal("Error asking for money :", err)
	}

	h1.WaitNewLedger()
	h1.WaitNewLedger()

	b, err := h1.GetBalance(util.KeyString(key2))
	if err != nil {
		t.Fatal("Error fetching balance :", err)
	}
	if b == 0 {
		t.Fatal("Money not transfered")
	}

}
