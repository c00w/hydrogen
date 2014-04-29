package libhydrogen

import (
	"crypto/sha512"
	"log"

	"libhydrogen/message"
	"libnode"
	"util"

	capnp "github.com/glycerine/go-capnproto"
)

type Handler interface {
	message.Verifier
	Handle(m message.Message)
}

type messagePasser struct {
	node        *libnode.Node
	handler     Handler
	newNeighbor chan *libnode.NeighborNode
	newMessage  chan message.Message
	neighbors   map[string]*libnode.NeighborNode
}

func newMessagePasser(n *libnode.Node, h Handler) *messagePasser {
	mp := &messagePasser{
		n,
		h,
		make(chan *libnode.NeighborNode),
		make(chan message.Message),
		make(map[string]*libnode.NeighborNode),
	}

	go mp.handleConns()
	go mp.handleMessages()
	n.AddListener("hydrogen", mp.newNeighbor)

	return mp
}

func (mp *messagePasser) handleConns() {
	for c := range mp.newNeighbor {
		mp.neighbors[c.Account()] = c
		go mp.handleConn(c)
	}
}

func (mp *messagePasser) handleConn(c *libnode.NeighborNode) {
	var seg *capnp.Segment
	var err error

	for {
		seg, err = capnp.ReadFromStream(c, nil)
		util.Debugf("Host %s, Message recieved", util.Short(util.KeyString(mp.node.Key)))
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
		mp.newMessage <- m
	}
}

func (mp *messagePasser) handleMessages() {
	for m := range mp.newMessage {
		go mp.handler.Handle(m)
		mp.passMessage(m)
	}
}

func (mp *messagePasser) SendChange(c message.Change) {
	m := message.CreateMessageFromChange(c, mp.node.Key)
	mp.newMessage <- m
}

func (mp *messagePasser) SendVote(v message.Vote) {
	m := message.CreateMessageFromVote(v, mp.node.Key)
	mp.newMessage <- m
}

func (mp *messagePasser) passMessage(m message.Message) {

	n := message.AppendAuthMessage(m, mp.node.Key)

	seen := make(map[string]bool)

	util.Debugf("Processing message %v", m)
	for _, a := range m.AuthChain().ToArray() {
		util.Debugf("Host %v seen", a)
		seen[a.Account()] = true
	}

	for name, neighbor := range mp.neighbors {
		if !seen[name] {
			util.Debugf("Host %v not seen, sending message", util.Short(name))
			n.WriteTo(neighbor)
		}
	}
}
