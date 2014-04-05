package libnode

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"testing"
)

func TestTLSConnection(t *testing.T) {
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	N := &Node{
		"account",
		priv,
		"ssl://test_machine:20",
        nil,
	}

	tc := make(chan *tls.Conn)
	N.TLSListen("127.0.0.1:2001", tc)

	c2 := N.TLSConnect("127.0.0.1:2001")
	c1 := <-tc

	n, err := c1.Write([]byte("Foo"))
	if err != nil {
		t.Fatal(err)
	}

	if n == 0 {
		t.Fatal("No data written")
	}

	b := make([]byte, 4)

	n, err = c2.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	if n == 0 {
		t.Fatal("No data read")
	}

	if b[0] != 'F' {
		t.Fatalf("Excepted \"Foo\", got %s", string(b))
	}

    _ = NewNeighborNode(c1)
}
