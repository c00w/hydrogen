package libhydrogen

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha512"
	"hash"
	"log"
	"time"

	"libhydrogen/message"
	"libnode"

	capnp "github.com/glycerine/go-capnproto"
)

type Hydrogen struct {
	node     *libnode.Node
	ledger   *Ledger
	tau      time.Duration
	key      *ecdsa.PrivateKey
	incoming chan *libnode.NeighborNode
	outgoing chan message.Message
}

func NewHydrogen(n *libnode.Node, key *ecdsa.PrivateKey,
	l *Ledger, tau time.Duration) *Hydrogen {
	h := &Hydrogen{
		n,
		l,
		tau,
		key,
		make(chan *libnode.NeighborNode),
		nil,
	}

	go h.handleConns()
	n.AddListener("hydrogen", h.incoming)

	return h

}

func (h *Hydrogen) handleConns() {
	for c := range h.incoming {
		go h.handleConn(c)
	}
}

func (h *Hydrogen) handleConn(c *libnode.NeighborNode) {
	buf := new(bytes.Buffer)
	var seg *capnp.Segment
	var err error

	for {
		seg, err = capnp.ReadFromStream(c, buf)
		if err != nil {
			panic(err)
		}
		m := message.ReadRootMessage(seg)
		s := sha512.New()
		if err := m.Verify(h.ledger, s); err != nil {
			log.Printf("Node %s: %s", h.node.Account, err)
			panic(err)
			continue
		}
		go h.HandleMessage(m, s)
		if h.outgoing != nil {
			h.outgoing <- m
		}
	}
}

func (h *Hydrogen) CreateMessageFromChange(c message.Change) *capnp.Segment {
	n := capnp.NewBuffer(nil)

	m := message.NewRootMessage(n)
	m.Payload().SetChange(c)

	run := sha512.New()
	c.Hash(run)

	a := message.NewSignedAuthorization(n, h.node.Account, h.key, run.Sum(nil))

	al := message.NewAuthorizationList(n, 1)
	capnp.PointerList(al).Set(0, capnp.Object(a))

	m.SetAuthChain(al)
	return n
}

func (h *Hydrogen) SendChange(c message.Change) {

	n := h.CreateMessageFromChange(c)

	for _, name := range h.node.ListNeighbors() {
		n.WriteTo(h.node.GetNeighbor(name))
	}
}

func (h *Hydrogen) AppendAuthMessage(m message.Message, run hash.Hash) (*capnp.Segment, message.Message) {

	n := capnp.NewBuffer(nil)

	m2 := message.NewRootMessage(n)

	l := message.NewAuthorizationList(n, m.AuthChain().Len()+1)
	for i, v := range m.AuthChain().ToArray() {
		capnp.PointerList(l).Set(i, capnp.Object(v))
	}

	a := message.NewSignedAuthorization(n, h.node.Account, h.key, run.Sum(nil))

	capnp.PointerList(l).Set(m.AuthChain().Len(), capnp.Object(a))

	m2.SetAuthChain(l)

	switch m.Payload().Which() {
	case message.MESSAGEPAYLOAD_VOTE:
		m2.Payload().SetVote(m.Payload().Vote())
	case message.MESSAGEPAYLOAD_CHANGE:
		m2.Payload().SetChange(m.Payload().Change())
	}

	return n, m2
}

func (h *Hydrogen) HandleMessage(m message.Message, run hash.Hash) {

	n, _ := h.AppendAuthMessage(m, run)

	seen := make(map[string]bool)

	for _, a := range m.AuthChain().ToArray() {
		seen[a.Account()] = true
	}

	for _, name := range h.node.ListNeighbors() {
		if !seen[name] {
			n.WriteTo(h.node.GetNeighbor(name))
		}
	}
}
