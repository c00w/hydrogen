package libnode

import (
	"crypto/ecdsa"
	"crypto/tls"
	"sync"

	"util"
)

type Node struct {
	Key      *ecdsa.PrivateKey
	Location string

	lock      *sync.RWMutex
	neighbors map[string]map[string]*NeighborNode
	listeners map[string]chan *NeighborNode
}

func NewNode(Key *ecdsa.PrivateKey, Location string) *Node {
	n := &Node{
		Key,
		Location,
		&sync.RWMutex{},
		make(map[string]map[string]*NeighborNode),
		make(map[string]chan *NeighborNode),
	}
	return n
}

func (n *Node) Listen(address string) {
	tc := make(chan *tls.Conn)
	n.tlsListen(address, tc)
	go n.handleConns(tc)
}

func (n *Node) Connect(address, protocol string) {
	c := n.tlsConnect(address, protocol)
	n.handleConn(c)
}

func (n *Node) handleConns(tc chan *tls.Conn) {
	for c := range tc {
		util.Debugf("Recieved connection")
		n.handleConn(c)
	}
}

func (n *Node) handleConn(c *tls.Conn) {
	N := NewNeighborNode(c)
	n.lock.Lock()
	defer n.lock.Unlock()
	if _, ok := n.neighbors[N.Protocol]; !ok {
		n.neighbors[N.Protocol] = make(map[string]*NeighborNode)
	}
	n.neighbors[N.Protocol][N.Account()] = N
	o, ok := n.listeners[N.Protocol]
	if ok {
		o <- N
	}
}

func (n *Node) protocols() []string {
	l := make([]string, 0, len(n.listeners))
	for p, _ := range n.listeners {
		l = append(l, p)
	}
	return l
}

func (n *Node) getNeighbor(account string, protocol string) *NeighborNode {
	n.lock.RLock()
	defer n.lock.RUnlock()
	return n.neighbors[protocol][account]
}

func (n *Node) listNeighbors() []string {
	n.lock.RLock()
	defer n.lock.RUnlock()
	nl := make([]string, 0, len(n.neighbors))
	for p, _ := range n.neighbors {
		for k, _ := range n.neighbors[p] {
			nl = append(nl, k)
		}
	}
	return nl
}

func (n *Node) AddListener(protocol string, c chan *NeighborNode) {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.listeners[protocol] = c
	for _, N := range n.neighbors[protocol] {
		c <- N
	}
}
