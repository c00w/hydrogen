package libnode

import (
	"crypto/ecdsa"
	"crypto/tls"
	"sync"
)

type Node struct {
	Account  string
	Key      *ecdsa.PrivateKey
	Location string

	lock      *sync.RWMutex
	neighbors map[string]*NeighborNode
}

func NewNode(Account string, Key *ecdsa.PrivateKey, Location string) *Node {
	n := &Node{
		Account,
		Key,
		Location,
		&sync.RWMutex{},
		make(map[string]*NeighborNode),
	}
	return n
}

func (n *Node) Listen(address string) {
	tc := make(chan *tls.Conn)
	n.tlsListen(address, tc)

	go n.handleconns(tc)
}

func (n *Node) Connect(address string) {
	c := n.tlsConnect(address)
	N := NewNeighborNode(c)
	n.lock.Lock()
	n.neighbors[N.Account] = N
	n.lock.Unlock()

}

func (n *Node) handleconns(tc chan *tls.Conn) {
	for c := range tc {
		N := NewNeighborNode(c)
		n.lock.Lock()
		n.neighbors[N.Account] = N
		n.lock.Unlock()
	}
}

func (n *Node) GetNeighbor(account string) *NeighborNode {
	n.lock.RLock()
	defer n.lock.RUnlock()
	return n.neighbors[account]
}

func (n *Node) ListNeighbors() []string {
	n.lock.RLock()
	defer n.lock.RUnlock()
	nl := make([]string, 0, len(n.neighbors))
	for k, _ := range n.neighbors {
		nl = append(nl, k)
	}
	return nl
}

func (n *Node) AddListener(protocol string, c chan *tls.Conn) {

}
