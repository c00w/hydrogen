package libhelium

import (
	"time"

	"libhydrogen"
	"libnode"
	"util"

	capnp "github.com/glycerine/go-capnproto"
)

func Connect(n *libnode.Node, address string) (*libhydrogen.Ledger, error) {
	tc := make(chan *libnode.NeighborNode)
	n.AddListener("helium", tc)
	go n.Connect(address, "helium")
	s := <-tc
	util.Debugf("Connected")

	ns, err := capnp.ReadFromStream(s, nil)
	if err != nil {
		return nil, err
	}
	util.Debugf("Segment Recieved")

	le := ReadRootLedger(ns)
	l := libhydrogen.NewLedger()
	l.Tau = time.Duration(le.Tau())
	l.Created = le.Created().Time()
	for _, v := range le.Accounts().ToArray() {
		l.AddEntry(string(v.Key()), v.Location(), v.Balance())
	}

	return l, nil
}
