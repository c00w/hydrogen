package libnode

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
)

func TestNode(t *testing.T) {
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	N1 := NewNode("account1", priv, "ssl://test_machine:20")

	N1.Listen("127.0.0.1:2002")

	N2 := NewNode("account2", priv, "ssl://test_machine:20")

	N2.Connect("127.0.0.1:2002")

	ns := N2.ListNeighbors()
	if len(ns) != 1 {
		t.Fatalf("Expected 1 Neighbor got %s", ns)
	}

	N2.GetNeighbor("account1").Write([]byte("Foo"))
	b := make([]byte, 4)
	N1.GetNeighbor("account2").Read(b)

	if b[0] != 'F' {
		t.Fatal("Expected leading F got %s", b)
	}

}
