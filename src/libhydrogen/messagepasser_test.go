package libhydrogen

import (
	"crypto/sha512"
	"testing"

	"libhydrogen/message"
	"libnode"
	"util"

	capnp "github.com/glycerine/go-capnproto"
)

type nullhandler struct{}

func (n nullhandler) Handle(m message.Message) {}

type channelhandler chan message.Message

func (c channelhandler) Handle(m message.Message) { c <- m }

func TestMessageManipulation(t *testing.T) {
	key1 := util.GenKey()
	key2 := util.GenKey()

	n1 := libnode.NewNode("node1", key1, "location1")
	n2 := libnode.NewNode("node2", key2, "location2")

	l := NewLedger()
	l.AddEntry("node1", util.KeyString(key1), "location1")
	l.AddEntry("node2", util.KeyString(key2), "location2")

	h1 := NewMessagePasser(n1, key1, l, nullhandler{})
	h2 := NewMessagePasser(n2, key2, l, nullhandler{})

	s1 := capnp.NewBuffer(nil)
	c := message.NewChange(s1)

	s2 := h1.CreateMessageFromChange(c)
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

	s3, m2o := h2.AppendAuthMessage(m, h)
	m2 := message.ReadRootMessage(s3)

	if m2.AuthChain().Len() != 2 {
		for i, v := range m.AuthChain().ToArray() {
			t.Logf("m.Authchain[%d] = %s", i, v.Account())
		}

		t.Logf("m20.len %d", m2o.AuthChain().Len())

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

	n1 := libnode.NewNode("node1", key1, "location1")
	n2 := libnode.NewNode("node2", key2, "location2")
	n3 := libnode.NewNode("node3", key3, "location3")

	l := NewLedger()
	l.AddEntry("node1", util.KeyString(key1), "location1")
	l.AddEntry("node2", util.KeyString(key2), "location2")
	l.AddEntry("node3", util.KeyString(key3), "location3")

	tc := make(chan message.Message)

	h1 := NewMessagePasser(n1, key1, l, nullhandler{})
	NewMessagePasser(n2, key2, l, nullhandler{})
	NewMessagePasser(n3, key3, l, channelhandler(tc))

	n1.Listen("127.0.0.1:3001")
	n2.Listen("127.0.0.1:3002")
	n2.Connect("127.0.0.1:3001")
	n3.Connect("127.0.0.1:3002")

	n := capnp.NewBuffer(nil)
	c := message.NewChange(n)

	h1.SendChange(c)
	m := <-tc
	if m.AuthChain().Len() != 2 {
		t.Fatal("message was not setup correctly")
	}
}