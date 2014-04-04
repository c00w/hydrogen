package libnode

import (
    "crypto/tls"
    "crypto/rand"
    "crypto/x509"
    "crypto/x509/pkix"
    "math/big"
    "net"
    "time"
)

import _ "crypto/sha512"

// Listen on an address and wrap incoming connections in TLS
func (n *Node) TLSListen(address string, t chan *tls.Conn) {
    rc := make(chan struct{})

    go func(n *Node, rc chan struct{}) {
        l, err := net.Listen("tcp", address)
        if err != nil {
            panic(err.Error())
        }
        rc <- struct{}{}

        for {
            c, err := l.Accept()
            if err != nil {
                panic(err.Error())
            }

            go func (n *Node, c net.Conn, t chan *tls.Conn) {
                nc := tls.Server(c, n.TLSConfig())
                err = nc.Handshake()

                if err != nil {
                    panic(err)
                }
                t <- nc
            }(n, c, t)
        }
    }(n, rc)
    <- rc
}

// Connect to an address and wrap incoming connections in TLS
func (n *Node) TLSConnect(address string) *tls.Conn {
    c, err := tls.Dial("tcp", address, n.TLSConfig())
    if err != nil {
        panic(err.Error())
    }
    return c
}

// Generate a tls config
func (n *Node) TLSConfig() *tls.Config {
    return &tls.Config{
        Certificates: []tls.Certificate{n.CreateTLSCert()},
        PreferServerCipherSuites: true,
        SessionTicketsDisabled: true,
        MinVersion: tls.VersionTLS12,
        ClientAuth: tls.RequireAnyClientCert,
        CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256},
        InsecureSkipVerify: true,
        NextProtos: []string{"hydrogen-core"},
    }
}

/*
This function generates TLS certs for a given key id. The cert will be self
Signed and the DNSNames will contain the location it is listening on. It will
also expire within 1 hour and should only be used for one connection
*/
func (n *Node) CreateTLSCert() tls.Certificate {

    NotBefore := time.Now().Add(-1 * time.Hour).UTC()
    NotAfter := time.Now().Add(time.Hour).UTC()

    template := &x509.Certificate{
        SerialNumber: new(big.Int).SetInt64(0),
        Subject: pkix.Name{CommonName: n.Account},
        NotBefore: NotBefore,
        NotAfter: NotAfter,
        KeyUsage:       x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
        ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
        BasicConstraintsValid: true,
        MaxPathLen: 1,
        IsCA: true,
        SubjectKeyId: []byte{1,2,3,4},
        Version: 2,
    }

    cert, err := x509.CreateCertificate(rand.Reader, template, template, &n.Key.PublicKey, n.Key)
    if err != nil {
        panic(err.Error())
    }

    tlscert := tls.Certificate{
        Certificate: [][]byte{cert},
        PrivateKey: n.Key,
    }
    return tlscert
}

