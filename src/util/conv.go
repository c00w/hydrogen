package util

func UInt64ToBA(u uint64) []byte {
	b := make([]byte, 8)
	for i, _ := range b {
		b[i] = byte((u >> uint(8*i)) & 0xff)
	}
	return b
}

func UInt32ToBA(u uint32) []byte {
	b := make([]byte, 4)
	for i, _ := range b {
		b[i] = byte((u >> uint(8*i)) & 0xff)
	}
	return b
}

func UInt16ToBA(u uint16) []byte {
	b := make([]byte, 2)
	for i, _ := range b {
		b[i] = byte((u >> uint(8*i)) & 0xff)
	}
	return b
}
