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
    tlscert := CreateTLSCert("account", "ssl://test_machine:20", &priv.PublicKey, priv)
    if tlscert == nil {
        t.Fatal(tlscert)
    }

}
