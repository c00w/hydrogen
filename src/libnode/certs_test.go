package libnode

import (
    "testing"
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
)

func TestTLSCertCreation(t *testing.T) {
    priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
    if err != nil {
        t.Fatal(err)
    }

    n := &Node{
        "account",
        priv,
        "ssl://test_machine:20",
    }

    tlscert := n.CreateTLSCert()
    tlscert = tlscert

}
