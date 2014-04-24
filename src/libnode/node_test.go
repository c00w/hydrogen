package libnode

import (
	"io"
	"testing"

	"util"
)

func TestNode(t *testing.T) {
	priv := util.GenKey()
	priv2 := util.GenKey()

	N1 := NewNode(priv, "ssl://test_machine:20")
	N1c := make(chan *NeighborNode)
	N1.AddListener("hydrogen", N1c)
	N1.Listen("127.0.0.1:2002")

	N2 := NewNode(priv2, "ssl://test_machine:20")
	N2c := make(chan *NeighborNode)
	N2.AddListener("hydrogen", N2c)
	go N2.Connect("127.0.0.1:2002", "hydrogen")

	<-N2c
	<-N1c

	ns := N2.listNeighbors()
	if len(ns) != 1 {
		t.Fatalf("Expected 1 Neighbor got %s", ns)
	}

	N2.getNeighbor(util.KeyString(priv), "hydrogen").Write([]byte("Foo"))
	b := make([]byte, 4)
	io.ReadAtLeast(N1.getNeighbor(util.KeyString(priv2), "hydrogen"), b, 1)

	if b[0] != 'F' {
		t.Fatal("Expected leading F got %s", b)
	}

}
