package blockchain

import (
	"fmt"
)

type Transaction struct {
	In  [SignatureSize]byte
	Out map[[SignatureSize]byte]float64
}

// DecodeTransaction from binary format
func DecodeTransaction(data []byte) (*Transaction, error) {
	return nil, fmt.Errorf("not implemented")
}
