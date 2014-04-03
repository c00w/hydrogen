package libnode

import (
    "crypto/rand"
    "crypto/tls"
    "crypto/x509"
    "crypto/x509/pkix"
    "math/big"
    "time"
)

import _ "crypto/sha512"

/*
This function generates TLS certs for a given key id. The cert will be self
Signed and the DNSNames will contain the location it is listening on. It will 
also expire within 1 hour and should only be used for one connection
*/
func CreateTLSCert(account, location string, pub, priv interface{}) *tls.Certificate {

    NotBefore := time.Now()
    NotAfter := time.Now().Add(24 * time.Hour)

    template := &x509.Certificate{
        SerialNumber: new(big.Int).SetBytes([]byte(account)),
        Subject: pkix.Name{CommonName:account},
        DNSNames: []string{location},
        NotBefore: NotBefore,
        NotAfter: NotAfter,
        KeyUsage:       x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
        ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
        BasicConstraintsValid: true,
    }

    cert, err := x509.CreateCertificate(rand.Reader, template, template, pub, priv)
    if err != nil {
        panic(err.Error())
    }

    tlscert := new(tls.Certificate)
    tlscert.Certificate = [][]byte{cert}
    return tlscert
}
