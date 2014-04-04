package libnode

import (
	"crypto/tls"
	"crypto/x509"
)

// A connection to a neighbor
type NeighborNode struct {
	c tls.Conn
}

func verify(c tls.Conn) error {
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

	// TODO: Verify account -> private key mapping in directory
	return nil

}

// Create a new Neighbor & verifies the tls connection
func NewNeighborNode(c tls.Conn) *NeighborNode {

	err := verify(c)
	if err != nil {
		panic(err.Error())
	}

	n := new(NeighborNode)
	n.c = c
	return n
}
