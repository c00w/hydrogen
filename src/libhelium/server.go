package libhelium

import (
	"libhydrogen"
	"libnode"
	"util"

	capnp "github.com/glycerine/go-capnproto"
)

type LedgerSource interface {
	WaitNewLedger() *libhydrogen.Ledger
}

type Server struct {
	node         *libnode.Node
	tc           chan *libnode.NeighborNode
	ledgerSource LedgerSource
}

func NewServer(n *libnode.Node, ls LedgerSource) *Server {
	tc := make(chan *libnode.NeighborNode)
	s := &Server{n, tc, ls}
	n.AddListener("helium", tc)
	go s.eventloop()
	return s
}

func (s *Server) eventloop() {
	for n := range s.tc {
		go s.dumpLedger(n)
	}
}

func (s *Server) dumpLedger(n *libnode.NeighborNode) {
	util.Debugf("Recieved Request")
	util.Debugf("Waiting for new ledger")
	l := s.ledgerSource.WaitNewLedger()

	util.Debugf("Encoding Ledger")
	b := capnp.NewBuffer(nil)
	le := NewRootLedger(b)
	le.SetTau(l.Tau.Nanoseconds())
	le.SetCreated(util.NewTimeFrom(b, l.Created))

	la := NewAccountList(b, len(l.Accounts))
	i := 0
	for _, v := range l.Accounts {
		a := NewAccount(b)
		a.SetKey([]byte(v.Key))
		a.SetLocation(v.Location)
		a.SetBalance(v.Balance)
		capnp.PointerList(la).Set(i, capnp.Object(a))
		i += 1
	}

	le.SetAccounts(la)
	util.Debugf("Writing ledger")
	_, err := b.WriteTo(n)
	if err != nil {
		util.Debugf("Error: %v", err)
	}
	util.Debugf("Done")
}
