package libhydrogen

import (
	"bytes"
	"crypto/ecdsa"
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
}

func NewHydrogen(n *libnode.Node, account string, key *ecdsa.PrivateKey,
	l *Ledger, tau time.Duration) *Hydrogen {
	h := &Hydrogen{
		n,
		l,
		tau,
		key,
		make(chan *libnode.NeighborNode),
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
			break
		}
		m := message.ReadRootMessage(seg)
		m = m
	}

	if err != nil {
		panic(err)
	}
}
