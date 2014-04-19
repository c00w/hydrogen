package libnode

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
)

// A connection to a neighbor
type NeighborNode struct {
	*tls.Conn
	Protocol string
}

func verifySigned(c *tls.Conn) error {
	certs := c.ConnectionState().PeerCertificates
	if len(certs) != 1 {
		panic("Weird certs")
	}

	cert := certs[0]

	certpool := x509.NewCertPool()
	certpool.AddCert(cert)

	_, err := cert.Verify(x509.VerifyOptions{
		Roots: certpool,
		KeyUsages: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		}})
	if err != nil {
		return err
	}

	return nil

}

// Create a new Neighbor & verifies the tls connection
func NewNeighborNode(c *tls.Conn) *NeighborNode {

	err := verifySigned(c)
	if err != nil {
		panic(err.Error())
	}

	n := new(NeighborNode)
	n.Conn = c
	n.Protocol = n.protocol()
	return n
}

func (n *NeighborNode) Account() string {
	s, _ := hex.DecodeString(n.Conn.ConnectionState().PeerCertificates[0].Subject.CommonName)
	return string(s)
}

func (n *NeighborNode) protocol() string {
	return n.Conn.ConnectionState().NegotiatedProtocol
}
