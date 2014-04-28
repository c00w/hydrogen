package message

import (
	"fmt"
	"util"
)

func shortb(b []byte) string {
	return util.Short(string(b))
}

func (m Message) String() string {
	return fmt.Sprintf("Message{%v, %v}", m.Payload(), m.AuthChain())
}

func (m MessagePayload) String() string {
	switch m.Which() {
	case MESSAGEPAYLOAD_VOTE:
		return fmt.Sprintf("Payload{%v, %v}", "VOTE", m.Vote())
	case MESSAGEPAYLOAD_CHANGE:
		return fmt.Sprintf("Payload{%v, %v}", "CHANGE", m.Change())
	default:
		return "Payload{UNKNOWN}"
	}
}

func (v Vote) String() string {
	return fmt.Sprintf("Vote{%v, %v, %v, %v, %v}", v.Votes(), v.Rate(), v.Drop(), v.Time(), v.Authorization())
}

func (dl DropChange_List) String() string {
	d := "Drop_List["
	for i := 0; i < dl.Len(); i++ {
		d += fmt.Sprintf("%v, ", dl.At(i))
	}
	d += "]"
	return d
}

func (d DropChange) String() string {
	return fmt.Sprintf("Drop{%s}", util.Short(string(d.Account())))
}

func (cl Change_List) String() string {
	s := "Change_List["
	for i := 0; i < cl.Len(); i++ {
		s += fmt.Sprintf("%v, ", cl.At(i))
	}
	s += "]"
	return s
}

func (a Authorization) String() string {
	return fmt.Sprintf("Authorization{%s...}", util.Short(a.Account()))
}

func (c Change) String() string {
	return fmt.Sprintf("Change{%v, %v, %v}", c.Type(), c.Authorization(), c.Created())
}

func (t ChangeType) String() string {
	switch t.Which() {
	case CHANGETYPE_TRANSACTION:
		return fmt.Sprintf("%v", t.Transaction())
	default:
		return fmt.Sprintf("UNKNOWN")
	}
}

func (t TransactionChange) String() string {
	return fmt.Sprintf("Transfer{%s..., %s..., %d}", shortb(t.Source()), shortb(t.Destination()), t.Amount())
}
