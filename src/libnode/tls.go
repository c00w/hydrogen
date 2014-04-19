package libnode

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"math/big"
	"net"
	"time"

	"util"
)

import _ "crypto/sha512"

// Listen on an address and wrap incoming connections in TLS
func (n *Node) tlsListen(address string, t chan *tls.Conn) {
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

			go func(n *Node, c net.Conn, t chan *tls.Conn) {
				nc := tls.Server(c, n.tlsConfig())
				err = nc.Handshake()

				if err != nil {
					panic(err)
				}
				t <- nc
			}(n, c, t)
		}
	}(n, rc)
	<-rc
}

// Connect to an address and wrap incoming connections in TLS
func (n *Node) tlsConnect(address string) *tls.Conn {
	c, err := tls.Dial("tcp", address, n.tlsConfig())
	if err != nil {
		panic(err.Error())
	}
	return c
}

// Generate a TLS config
func (n *Node) tlsConfig() *tls.Config {
	return &tls.Config{
		Certificates:             []tls.Certificate{n.tlsCert()},
		PreferServerCipherSuites: true,
		SessionTicketsDisabled:   true,
		MinVersion:               tls.VersionTLS12,
		ClientAuth:               tls.RequireAnyClientCert,
		CipherSuites:             []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256},
		InsecureSkipVerify:       true,
		NextProtos:               []string{"hydrogen"},
	}
}

/*
This function generates TLS certs for a given key id. The cert will be self
Signed and the DNSNames will contain the location it is listening on. It will
also expire within 1 hour and should only be used for one connection
*/
func (n *Node) tlsCert() tls.Certificate {

	NotBefore := time.Now().Add(-1 * time.Hour).UTC()
	NotAfter := time.Now().Add(time.Hour).UTC()

	template := &x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject:      pkix.Name{CommonName: hex.EncodeToString([]byte(util.KeyString(n.Key)))},
		NotBefore:    NotBefore,
		NotAfter:     NotAfter,
		KeyUsage:     x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		MaxPathLen:   1,
		IsCA:         true,
		BasicConstraintsValid: true,
	}

	cert, err := x509.CreateCertificate(rand.Reader, template, template, &n.Key.PublicKey, n.Key)
	if err != nil {
		panic(err.Error())
	}

	tlscert := tls.Certificate{
		Certificate: [][]byte{cert},
		PrivateKey:  n.Key,
	}
	return tlscert
}
