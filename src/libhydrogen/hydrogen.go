package libhydrogen

import (
	"bytes"
	"crypto/ecdsa"
    "crypto/sha512"
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
        s = sha512.New()
        if !m.Verify(h.ledger, s) {
            continue
        }
        h.HandleMessage(m, s)
	}

	if err != nil {
		panic(err)
	}
}

func (h *Hydrogen) HandleMessage(m message.Message, s hash.Hash) {

    n := capnp.NewBuffer(nil)
    
    l := message.NewAuthorizationList(n, m.Authchain().Len()+1)
    for i, v := range(message.Authchain().ToArray()) {
        capnp.PointerList(l).Set(i, capnp.Object(v))
    }

    key := NewKey(n)
    key.SetX(h.key.X.Bytes())
    key.SetY(h.key.Y.Bytes())

    r, s, err := ecdsa.Sign(rand, h.key, s.Sum())

    sig := message.NewSignature(n)
    sig.SetR(r.Bytes())
    sig.SetS(r.Bytes())

    keysig := Message.NewKeySignature(n)
    keysig.SetKey(key)
    keysig.SetSignature(sig)

    sl := message.NewKeySignatureList(n, 1)
    capnp.PointerList(sl).Set(0, capnp.Object(keysig))

    a := NewAuthorization(n)
    a.SetAccount(h.Account)
    a.SetKeySignatureList(l)

    capnp.PointerList(l).Set(0, capnp.Object(a))

    m2 := message.NewRootMessage(n)
    switch m.Payload().Which() {
        case message.MESSAGEPAYLOAD_VOTE:
            m2.Payload().SetVote(m.Payload().Vote())
        case message.MESSAGEPAYLOAD_CHANGE:
            m2.Payload().SetChane(m.Payload().Vote())
    }
    m2.SetAuthChain(l)

    s := make(map[string]bool)

    for _, a := range(m.AuthChain().ToArray()) {
        s.[a.Account()] = true
    }

    for _, h := h.Node.ListNeighbors() {
        if !s[h.Account] {
            m2.WriteTo(h)
        }
    }
}
