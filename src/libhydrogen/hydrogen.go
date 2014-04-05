package libhydrogen

import (
	"crypto/ecdsa"
	"crypto/tls"
	"libnode"
	"time"
)

type Hydrogen struct {
	node     *libnode.Node
	ledger   *Ledger
	tau      time.Duration
	key      *ecdsa.PrivateKey
	incoming chan *tls.Conn
}

func NewHydrogen(n *libnode.Node, account string, key *ecdsa.PrivateKey,
	l *Ledger, tau time.Duration) *Hydrogen {
	h := &Hydrogen{
		n,
		l,
		tau,
		key,
		make(chan *tls.Conn),
	}

	n.AddListener("hydrogen", h.incoming)

	go h.handleConns()
	return h

}

func (h *Hydrogen) handleConns() {
	for c := range h.incoming {
		go h.handleConn(c)
	}
}

func (h *Hydrogen) handleConn(c *tls.Conn) {
}
