package libnode

// AUTO GENERATED - DO NOT EDIT

import (
	"bufio"
	"bytes"
	"encoding/json"
	C "github.com/jmckaskill/go-capnproto"
	"io"
	"unsafe"
)

type Signature C.Struct

func NewSignature(s *C.Segment) Signature      { return Signature(s.NewStruct(0, 2)) }
func NewRootSignature(s *C.Segment) Signature  { return Signature(s.NewRootStruct(0, 2)) }
func ReadRootSignature(s *C.Segment) Signature { return Signature(s.Root(0).ToStruct()) }
func (s Signature) Key() []byte                { return C.Struct(s).GetObject(0).ToData() }
func (s Signature) SetKey(v []byte)            { C.Struct(s).SetObject(0, s.Segment.NewData(v)) }
func (s Signature) Signature() []byte          { return C.Struct(s).GetObject(1).ToData() }
func (s Signature) SetSignature(v []byte)      { C.Struct(s).SetObject(1, s.Segment.NewData(v)) }
func (s Signature) WriteJSON(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	err = b.WriteByte('{')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"key\":")
	if err != nil {
		return err
	}
	{
		s := s.Key()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"signature\":")
	if err != nil {
		return err
	}
	{
		s := s.Signature()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte('}')
	if err != nil {
		return err
	}
	err = b.Flush()
	return err
}
func (s Signature) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteJSON(&b)
	return b.Bytes(), err
}

type Signature_List C.PointerList

func NewSignatureList(s *C.Segment, sz int) Signature_List {
	return Signature_List(s.NewCompositeList(0, 2, sz))
}
func (s Signature_List) Len() int           { return C.PointerList(s).Len() }
func (s Signature_List) At(i int) Signature { return Signature(C.PointerList(s).At(i).ToStruct()) }
func (s Signature_List) ToArray() []Signature {
	return *(*[]Signature)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type Time C.Struct

func NewTime(s *C.Segment) Time        { return Time(s.NewStruct(16, 0)) }
func NewRootTime(s *C.Segment) Time    { return Time(s.NewRootStruct(16, 0)) }
func ReadRootTime(s *C.Segment) Time   { return Time(s.Root(0).ToStruct()) }
func (s Time) Seconds() uint64         { return C.Struct(s).Get64(0) }
func (s Time) SetSeconds(v uint64)     { C.Struct(s).Set64(0, v) }
func (s Time) NanoSeconds() uint32     { return C.Struct(s).Get32(8) }
func (s Time) SetNanoSeconds(v uint32) { C.Struct(s).Set32(8, v) }
func (s Time) WriteJSON(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	err = b.WriteByte('{')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"seconds\":")
	if err != nil {
		return err
	}
	{
		s := s.Seconds()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"nanoSeconds\":")
	if err != nil {
		return err
	}
	{
		s := s.NanoSeconds()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte('}')
	if err != nil {
		return err
	}
	err = b.Flush()
	return err
}
func (s Time) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteJSON(&b)
	return b.Bytes(), err
}

type Time_List C.PointerList

func NewTimeList(s *C.Segment, sz int) Time_List { return Time_List(s.NewCompositeList(16, 0, sz)) }
func (s Time_List) Len() int                     { return C.PointerList(s).Len() }
func (s Time_List) At(i int) Time                { return Time(C.PointerList(s).At(i).ToStruct()) }
func (s Time_List) ToArray() []Time              { return *(*[]Time)(unsafe.Pointer(C.PointerList(s).ToArray())) }

type Change C.Struct

func NewChange(s *C.Segment) Change                { return Change(s.NewStruct(8, 5)) }
func NewRootChange(s *C.Segment) Change            { return Change(s.NewRootStruct(8, 5)) }
func ReadRootChange(s *C.Segment) Change           { return Change(s.Root(0).ToStruct()) }
func (s Change) Type() uint8                       { return C.Struct(s).Get8(0) }
func (s Change) SetType(v uint8)                   { C.Struct(s).Set8(0, v) }
func (s Change) Account() []byte                   { return C.Struct(s).GetObject(0).ToData() }
func (s Change) SetAccount(v []byte)               { C.Struct(s).SetObject(0, s.Segment.NewData(v)) }
func (s Change) Authorization() Signature_List     { return Signature_List(C.Struct(s).GetObject(1)) }
func (s Change) SetAuthorization(v Signature_List) { C.Struct(s).SetObject(1, C.Object(v)) }
func (s Change) Created() Time                     { return Time(C.Struct(s).GetObject(2).ToStruct()) }
func (s Change) SetCreated(v Time)                 { C.Struct(s).SetObject(2, C.Object(v)) }
func (s Change) NewValue() []byte                  { return C.Struct(s).GetObject(3).ToData() }
func (s Change) SetNewValue(v []byte)              { C.Struct(s).SetObject(3, s.Segment.NewData(v)) }
func (s Change) Destination() []byte               { return C.Struct(s).GetObject(4).ToData() }
func (s Change) SetDestination(v []byte)           { C.Struct(s).SetObject(4, s.Segment.NewData(v)) }
func (s Change) WriteJSON(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	err = b.WriteByte('{')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"type\":")
	if err != nil {
		return err
	}
	{
		s := s.Type()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"account\":")
	if err != nil {
		return err
	}
	{
		s := s.Account()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"authorization\":")
	if err != nil {
		return err
	}
	{
		s := s.Authorization()
		{
			err = b.WriteByte('[')
			if err != nil {
				return err
			}
			for i, s := range s.ToArray() {
				if i != 0 {
					_, err = b.WriteString(", ")
				}
				if err != nil {
					return err
				}
				err = s.WriteJSON(b)
				if err != nil {
					return err
				}
			}
			err = b.WriteByte(']')
		}
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"created\":")
	if err != nil {
		return err
	}
	{
		s := s.Created()
		err = s.WriteJSON(b)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"newValue\":")
	if err != nil {
		return err
	}
	{
		s := s.NewValue()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"destination\":")
	if err != nil {
		return err
	}
	{
		s := s.Destination()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte('}')
	if err != nil {
		return err
	}
	err = b.Flush()
	return err
}
func (s Change) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteJSON(&b)
	return b.Bytes(), err
}

type Change_List C.PointerList

func NewChangeList(s *C.Segment, sz int) Change_List { return Change_List(s.NewCompositeList(8, 5, sz)) }
func (s Change_List) Len() int                       { return C.PointerList(s).Len() }
func (s Change_List) At(i int) Change                { return Change(C.PointerList(s).At(i).ToStruct()) }
func (s Change_List) ToArray() []Change {
	return *(*[]Change)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type Message C.Struct

func NewMessage(s *C.Segment) Message           { return Message(s.NewStruct(0, 3)) }
func NewRootMessage(s *C.Segment) Message       { return Message(s.NewRootStruct(0, 3)) }
func ReadRootMessage(s *C.Segment) Message      { return Message(s.Root(0).ToStruct()) }
func (s Message) Votes() Change_List            { return Change_List(C.Struct(s).GetObject(0)) }
func (s Message) SetVotes(v Change_List)        { C.Struct(s).SetObject(0, C.Object(v)) }
func (s Message) Time() Time                    { return Time(C.Struct(s).GetObject(1).ToStruct()) }
func (s Message) SetTime(v Time)                { C.Struct(s).SetObject(1, C.Object(v)) }
func (s Message) Signature() Signature_List     { return Signature_List(C.Struct(s).GetObject(2)) }
func (s Message) SetSignature(v Signature_List) { C.Struct(s).SetObject(2, C.Object(v)) }
func (s Message) WriteJSON(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	err = b.WriteByte('{')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"votes\":")
	if err != nil {
		return err
	}
	{
		s := s.Votes()
		{
			err = b.WriteByte('[')
			if err != nil {
				return err
			}
			for i, s := range s.ToArray() {
				if i != 0 {
					_, err = b.WriteString(", ")
				}
				if err != nil {
					return err
				}
				err = s.WriteJSON(b)
				if err != nil {
					return err
				}
			}
			err = b.WriteByte(']')
		}
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"time\":")
	if err != nil {
		return err
	}
	{
		s := s.Time()
		err = s.WriteJSON(b)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"signature\":")
	if err != nil {
		return err
	}
	{
		s := s.Signature()
		{
			err = b.WriteByte('[')
			if err != nil {
				return err
			}
			for i, s := range s.ToArray() {
				if i != 0 {
					_, err = b.WriteString(", ")
				}
				if err != nil {
					return err
				}
				err = s.WriteJSON(b)
				if err != nil {
					return err
				}
			}
			err = b.WriteByte(']')
		}
		if err != nil {
			return err
		}
	}
	err = b.WriteByte('}')
	if err != nil {
		return err
	}
	err = b.Flush()
	return err
}
func (s Message) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteJSON(&b)
	return b.Bytes(), err
}

type Message_List C.PointerList

func NewMessageList(s *C.Segment, sz int) Message_List {
	return Message_List(s.NewCompositeList(0, 3, sz))
}
func (s Message_List) Len() int         { return C.PointerList(s).Len() }
func (s Message_List) At(i int) Message { return Message(C.PointerList(s).At(i).ToStruct()) }
func (s Message_List) ToArray() []Message {
	return *(*[]Message)(unsafe.Pointer(C.PointerList(s).ToArray()))
}
