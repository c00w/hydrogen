package liblithium

import (
	"net"

	capnp "github.com/glycerine/go-capnproto"
)

type Client struct {
	fd net.Conn
}

func NewClient() (*Client, error) {
	return NewClientAt("/run/hydrogend/lithium.socket")
}

func NewClientAt(location string) (*Client, error) {
	fd, err := net.Dial("unix", location)
	if err != nil {
		return nil, err
	}
	c := &Client{fd}
	return c, nil
}

func (c *Client) GetBalance(host string) string {
	b := capnp.NewBuffer(nil)
	com := NewRootCommand(b)
	g := NewGetBalance(b)
	g.SetAccount([]byte(host))
	com.SetGetbalance(g)
	b.WriteTo(c.fd)

	seg, err := capnp.ReadFromStream(c.fd, nil)
	if err != nil {
		return string(err.Error())
	}
	res := ReadRootResult(seg)
	return string(res.Message())
}

func (c *Client) SendMoney(to string, balance uint64) string {
	b := capnp.NewBuffer(nil)
	com := NewRootCommand(b)
	s := NewSendMoney(b)
	s.SetTo([]byte(to))
	s.SetAmount(balance)
	com.SetSendmoney(s)
	b.WriteTo(c.fd)

	seg, err := capnp.ReadFromStream(c.fd, nil)
	if err != nil {
		return string(err.Error())
	}
	res := ReadRootResult(seg)
	return string(res.Message())
}
