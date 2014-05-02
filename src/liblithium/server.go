package liblithium

import (
	"fmt"
	"log"
	"net"

	"libhydrogen"
	"util"

	capnp "github.com/glycerine/go-capnproto"
)

type Server struct {
	h    *libhydrogen.Hydrogen
	sock net.Listener
}

func NewServer(h *libhydrogen.Hydrogen) (*Server, error) {
	return NewServerAt(h, "/run/hydrogend/lithium.socket")
}

func NewServerAt(h *libhydrogen.Hydrogen, location string) (*Server, error) {
	fd, err := net.Listen("unix", location)
	if err != nil {
		return nil, err
	}
	s := &Server{h, fd}
	go s.eventloop()
	return s, nil
}

func (s *Server) eventloop() {
	for {
		c, err := s.sock.Accept()
		if err != nil {
			log.Println("Error accepting connection: %v", err)
			return
		}
		go s.handleConn(c)
	}
}

func (s *Server) handleConn(c net.Conn) {
	for {
		seg, err := capnp.ReadFromStream(c, nil)
		if err != nil {
			log.Print(err)
			return
		}
		command := ReadRootCommand(seg)
		r := ""
		util.Debugf("Command Recieved")
		switch command.Which() {
		case COMMAND_GETBALANCE:
			util.Debugf("Getbalance")
			getb := command.Getbalance()
			account := string(getb.Account())
			if len(account) == 0 {
				account = s.h.Account()
			}
			b, err := s.h.GetBalance(account)
			if err != nil {
				r += fmt.Sprintf("Error fetching balance: %v", err)
			} else {
				r += fmt.Sprintf("%d", b)
			}
		case COMMAND_SENDMONEY:
			util.Debugf("SendMoney")
			sm := command.Sendmoney()
			to := string(sm.To())
			amount := sm.Amount()
			err := s.h.TransferMoney(to, amount)
			if err != nil {
				r += fmt.Sprintf("Error sending money: %v", err)
			} else {
				r += "OK"
			}
		}

		util.Debugf("Response is %s", r)

		s := capnp.NewBuffer(nil)
		res := NewRootResult(s)
		res.SetMessage(r)
		s.WriteTo(c)
	}
}
