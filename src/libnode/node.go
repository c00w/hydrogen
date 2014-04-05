package libnode

import (
	"crypto/ecdsa"
	"crypto/tls"
)

type Node struct {
	Account  string
	Key      *ecdsa.PrivateKey
	Location string

	Neighbors map[string]*NeighborNode
}

func (n *Node) Listen(address string) {
	tc := make(chan *tls.Conn)
	n.tlsListen(address, tc)

	go n.handleconns(tc)
}

func (n *Node) Connect(address string) {
	c := n.tlsConnect(address)
	N := NewNeighborNode(c)
	n.Neighbors[N.Account] = N

}

func (n *Node) handleconns(tc chan *tls.Conn) {
	for c := range tc {
		N := NewNeighborNode(c)
		n.Neighbors[N.Account] = N
	}
}
