package libhydrogen

import (
	"crypto/sha512"
	"testing"

	"libhydrogen/message"
	"libnode"
	"util"

	capnp "github.com/glycerine/go-capnproto"
)

type nullhandler struct {
	*Ledger
}

func (n nullhandler) Handle(m message.Message)      {}
func (n nullhandler) RegisterBus(mp *MessagePasser) {}

type channelhandler struct {
	*Ledger
	c chan message.Message
}

func (c channelhandler) Handle(m message.Message)      { c.c <- m }
func (n channelhandler) RegisterBus(mp *MessagePasser) {}

func TestMessageManipulation(t *testing.T) {
	key1 := util.GenKey()
	key2 := util.GenKey()

	l := NewLedger()
	l.AddEntry(util.KeyString(key1), "location1", 100)
	l.AddEntry(util.KeyString(key2), "location2", 100)

	s1 := capnp.NewBuffer(nil)
	c := message.NewChange(s1)

	s2 := message.CreateMessageFromChange(c, key1)
	m := message.ReadRootMessage(s2)

	err := m.Verify(l, sha512.New())
	if err != nil {
		t.Fatalf("Verifying failed: %s", err.Error())
	}

	if m.AuthChain().Len() != 1 {
		t.Fatalf("len(m.Authchain) %d != 1", m.AuthChain().Len())
	}

	h := sha512.New()
	c.Hash(h)
	m.AuthChain().ToArray()[0].Hash(h)

	s3 := message.AppendAuthMessage(m, h, key2)
	m2 := message.ReadRootMessage(s3)

	if m2.AuthChain().Len() != 2 {
		for i, v := range m.AuthChain().ToArray() {
			t.Logf("m.Authchain[%d] = %s", i, v.Account())
		}

		for i, v := range m2.AuthChain().ToArray() {
			t.Logf("m2.Authchain[%d] = %s", i, v.Account())
		}
		t.Fatalf("len(m2.Authchain) %d != 2", m2.AuthChain().Len())
	}
}

func TestMessagePassing(t *testing.T) {
	key1 := util.GenKey()
	key2 := util.GenKey()
	key3 := util.GenKey()

	n1 := libnode.NewNode(key1, "location1")
	n2 := libnode.NewNode(key2, "location2")
	n3 := libnode.NewNode(key3, "location3")

	l := NewLedger()
	l.AddEntry(util.KeyString(key1), "location1", 100)
	l.AddEntry(util.KeyString(key2), "location2", 100)
	l.AddEntry(util.KeyString(key3), "location3", 100)

	tc := make(chan message.Message)

	h1 := NewMessagePasser(n1, nullhandler{l})
	NewMessagePasser(n2, nullhandler{l})
	NewMessagePasser(n3, channelhandler{l, tc})

	n1.Listen("127.0.0.1:3001")
	n2.Listen("127.0.0.1:3002")
	n2.Connect("127.0.0.1:3001", "hydrogen")
	n3.Connect("127.0.0.1:3002", "hydrogen")

	n := capnp.NewBuffer(nil)
	c := message.NewChange(n)

	h1.SendChange(c)
	m := <-tc
	if m.AuthChain().Len() != 2 {
		t.Errorf("Authchain incorrect length, %d != 2", m.AuthChain().Len())
		t.Errorf("Message was %v", m)
	}

	n = capnp.NewBuffer(nil)
	v := message.NewVote(n)
	h1.SendVote(v)
	m = <-tc
	if m.AuthChain().Len() != 2 {
		t.Errorf("Authchain incorrect length, %d != 2", m.AuthChain().Len())
		t.Errorf("Message was %v", m)
	}
}
