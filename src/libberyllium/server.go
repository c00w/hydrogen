package libberyllium

import (
	"libhydrogen"
	"libnode"
	"util"

	capnp "github.com/glycerine/go-capnproto"
)

type Server struct {
	n  *libnode.Node
	tc chan *libnode.NeighborNode
	h  *libhydrogen.Hydrogen
}

func NewServer(n *libnode.Node, h *libhydrogen.Hydrogen) *Server {
	tc := make(chan *libnode.NeighborNode)
	n.AddListener("beryllium", tc)
	s := &Server{n, tc, h}
	go s.eventloop()
	return s
}

func (s *Server) eventloop() {
	for c := range s.tc {
		go s.donateMoney(c)
	}
}

func (s *Server) donateMoney(c *libnode.NeighborNode) {
	util.Debugf("Client connected")
	ns, err := capnp.ReadFromStream(c, nil)
	if err != nil {
		util.Debugf("%v", err)
	}

	req := ReadRootRequest(ns)
	acc := string(req.Account())

	util.Debugf("Transferring money")
	s.h.TransferMoney(acc, 100)
	util.Debugf("Transfered")

	b := capnp.NewBuffer(nil)

	r := NewRootResponse(b)
	r.SetOk(true)
	b.WriteTo(c)
	util.Debugf("Sent")
}
