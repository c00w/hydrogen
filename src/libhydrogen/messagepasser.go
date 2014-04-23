package libhydrogen

import (
	"crypto/sha512"
	"hash"
	"log"

	"libhydrogen/message"
	"libnode"
	"util"

	capnp "github.com/glycerine/go-capnproto"
)

type Handler interface {
	message.Verifier
	Handle(m message.Message)
	RegisterBus(mp *messagePasser)
}

type messagePasser struct {
	node        *libnode.Node
	handler     Handler
	newNeighbor chan *libnode.NeighborNode
	newMessage  chan message.Message
}

func newMessagePasser(n *libnode.Node, h Handler) *messagePasser {
	mp := &messagePasser{
		n,
		h,
		make(chan *libnode.NeighborNode),
		make(chan message.Message),
	}

	go mp.handleConns()
	go mp.handleMessages()
	n.AddListener("hydrogen", mp.newNeighbor)
	h.RegisterBus(mp)

	return mp

}

func (mp *messagePasser) handleConns() {
	for c := range mp.newNeighbor {
		go mp.handleConn(c)
	}
}

func (mp *messagePasser) handleConn(c *libnode.NeighborNode) {
	var seg *capnp.Segment
	var err error

	for {
		seg, err = capnp.ReadFromStream(c, nil)
		if err != nil {
			panic(err)
		}
		m := message.ReadRootMessage(seg)
		s := sha512.New()
		if err := m.Verify(mp.handler, s); err != nil {
			log.Printf("Node %s: %s", util.KeyString(mp.node.Key), err)
			panic(err)
			continue
		}
		go mp.passMessage(m, s)
		mp.newMessage <- m
	}
}

func (mp *messagePasser) handleMessages() {
	for m := range mp.newMessage {
		mp.handler.Handle(m)
	}
}

func (mp *messagePasser) SendChange(c message.Change) {
	n := message.CreateMessageFromChange(c, mp.node.Key)
	mp.sendMessage(n)
}

func (mp *messagePasser) SendVote(v message.Vote) {
	n := message.CreateMessageFromVote(v, mp.node.Key)
	mp.sendMessage(n)
}

func (mp *messagePasser) sendMessage(n *capnp.Segment) {
	for _, name := range mp.node.ListNeighbors() {
		n.WriteTo(mp.node.GetNeighbor(name))
	}
	mp.newMessage <- message.ReadRootMessage(n)
}

func (mp *messagePasser) passMessage(m message.Message, run hash.Hash) {

	n := message.AppendAuthMessage(m, run, mp.node.Key)

	seen := make(map[string]bool)

	util.Debugf("Processing message %v", m)
	for _, a := range m.AuthChain().ToArray() {
		util.Debugf("Host %v seen", a)
		seen[a.Account()] = true
	}

	for _, name := range mp.node.ListNeighbors() {
		if !seen[name] {
			util.Debugf("Host %v not seen, sending message", util.Short(name))
			n.WriteTo(mp.node.GetNeighbor(name))
		}
	}
}
