package libnode

import (
	"crypto/ecdsa"
	"crypto/tls"
	"sync"

	"util"
)

/*
A local node listening on a socket for tls connections for specific protocols
from neighboring nodes.
*/
type Node struct {
	Key      *ecdsa.PrivateKey
	Location string

	lock      *sync.RWMutex
	neighbors map[string]map[string]*NeighborNode
	listeners map[string]chan *NeighborNode
}

// Create a newnode for a given account and location
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

// Listen on a given address for tls connections
func (n *Node) Listen(address string) {
	tc := make(chan *tls.Conn)
	n.tlsListen(address, tc)
	go n.handleConns(tc)
}

// Connect to a node at the given address with the following protocol
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

/*
Register a handler for a given protocol. All nieghbornodes will be sent to the
channel which support a given protocol, this includes already existing connections
on that protocol
*/
func (n *Node) AddListener(protocol string, c chan *NeighborNode) {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.listeners[protocol] = c
	for _, N := range n.neighbors[protocol] {
		c <- N
	}
}
