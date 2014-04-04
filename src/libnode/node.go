package libnode

import (
	"crypto/ecdsa"
)

type Node struct {
	Account  string
	Key      *ecdsa.PrivateKey
	Location string
}
