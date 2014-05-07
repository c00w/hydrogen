package libberyllium

import (
	"errors"

	"libnode"
	"util"

	capnp "github.com/glycerine/go-capnproto"
)

func GetMoney(n *libnode.Node, addr, account string) error {
	tc := make(chan *libnode.NeighborNode)
	n.AddListener("beryllium", tc)
	go n.Connect(addr, "beryllium")
	s := <-tc
	util.Debugf("Connected")

	b := capnp.NewBuffer(nil)
	req := NewRootRequest(b)
	req.SetAccount([]byte(account))
	b.WriteTo(s)
	util.Debugf("Request sent")

	b2, err := capnp.ReadFromStream(s, nil)

	util.Debugf("Response recieved")

	if err != nil {
		return err
	}
	resp := ReadRootResponse(b2)
	if !resp.Ok() {
		return errors.New("Unable to beg for money")
	}

	return nil
}
