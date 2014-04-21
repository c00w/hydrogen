package libhelium

import (
	"testing"

	"libhydrogen"
	"libnode"
	"util"
)

type DummyLedger struct {
	*libhydrogen.Ledger
}

func (d DummyLedger) GetLedger() *libhydrogen.Ledger {
	return d.Ledger
}

func TestDownload(t *testing.T) {

	key1 := util.GenKey()
	key2 := util.GenKey()

	l := libhydrogen.NewLedger()
	l.AddEntry(util.KeyString(key1), "127.0.0.1:3010", 100)
	ls := DummyLedger{l}

	n1 := libnode.NewNode(key1, "127.0.0.1:3010")
	n1.Listen("127.0.0.1:3010")
	n2 := libnode.NewNode(key2, "127.0.0.1:3011")

	NewServer(n1, ls)

	lr := Connect(n2, "127.0.0.1:3010")
	if len(lr.Accounts) != 1 {
		t.Fatalf("Account length incorrect, %v", lr.Accounts)
	}
}
