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
	listeners map[string]chan *NeighborNode
}

func NewNode(Account string, Key *ecdsa.PrivateKey, Location string) *Node {
	n := &Node{
		Account,
		Key,
		Location,
		&sync.RWMutex{},
		make(map[string]*NeighborNode),
		make(map[string]chan *NeighborNode),
	}
	return n
}

func (n *Node) Listen(address string) {
	tc := make(chan *tls.Conn)
	n.tlsListen(address, tc)
	go n.handleConns(tc)
}

func (n *Node) Connect(address string) {
	c := n.tlsConnect(address)
	n.handleConn(c)
}

func (n *Node) handleConns(tc chan *tls.Conn) {
	for c := range tc {
		n.handleConn(c)
	}
}

func (n *Node) handleConn(c *tls.Conn) {
	N := NewNeighborNode(c)
	n.lock.Lock()
	defer n.lock.Unlock()
	n.neighbors[N.Account] = N
	o, ok := n.listeners[N.Protocol]
	if ok {
		o <- N
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

func (n *Node) AddListener(protocol string, c chan *NeighborNode) {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.listeners[protocol] = c
	for _, N := range n.neighbors {
		if N.Protocol == protocol {
			c <- N
		}
	}
}
