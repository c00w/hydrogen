package libnode

import (
    "crypto/tls"
    "net"
)

// Listen on an address and wrap incoming connections in TLS
func (n *Node) TLSListen(address string) {
    l, err := net.Listen("tcp", address)
    if err != nil {
        panic(err.Error())
    }

    for {
        c, err := l.Accept()
        if err != nil {
            panic(err.Error())
        }

        go n.TLSSetup(c, true)
    }
}

// Connect to an address and wrap incoming connections in TLS
func (n *Node) TLSConnect(address string) {
    c, err := net.Dial("tcp", address)
    if err != nil {
        panic(err.Error())
    }

    go n.TLSSetup(c, false)
}

// Wrap a network connection in TLS
func (n *Node) TLSSetup(c net.Conn, server bool) *tls.Conn {

    config := &tls.Config{
        Certificates: []tls.Certificate{n.CreateTLSCert()},
        PreferServerCipherSuites: true,
        SessionTicketsDisabled: true,
        ServerName: n.Account,
        MinVersion: tls.VersionTLS12,
        ClientAuth: tls.RequireAnyClientCert,
        CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256},
        InsecureSkipVerify: true,
    }

    if server {
        return tls.Server(c, config)
    } else {
        return tls.Client(c, config)
    }
}
