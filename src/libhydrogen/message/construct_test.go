package message

import (
	"testing"

	"util"
)

func TestTransferStringify(t *testing.T) {
	key := util.GenKey()

	NewSignedTransaction(key, "foo", 10).String()
}
