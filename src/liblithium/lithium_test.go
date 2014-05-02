package liblithium

import (
	"os"
	"testing"

	"libhydrogen"
	"libnode"
	"util"
)

func TestCommands(t *testing.T) {
	key := util.GenKey()

	n := libnode.NewNode(key, "location1")

	l := libhydrogen.NewLedger()
	l.AddEntry(util.KeyString(key), "location1", 100)

	h := libhydrogen.NewHydrogen(n)
	h.SetBootStrap()
	h.AddLedger(l)

	sock := os.TempDir() + "/lithium.socket"
	_, err := NewServerAt(h, sock)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(sock)

	client, err := NewClientAt(sock)
	if err != nil {
		t.Fatal(err)
	}

	bal := client.GetBalance("")
	if bal != "100" {
		t.Errorf("GetBalance response is '%s' != 100", bal)
	}

	send := client.SendMoney("foo", 10)
	if send != "OK" {
		t.Errorf("SendMoney Response is '%s' != OK", send)
	}

	h.WaitNewLedger()
	h.WaitNewLedger()

	bal = client.GetBalance("")
	if bal != "91" {
		t.Errorf("GetBalance response is '%s' != 91", bal)
	}

	newb := client.GetBalance("foo")
	if newb != "9" {
		t.Errorf("GetBalance response is '%s' != 9", newb)
	}

}
