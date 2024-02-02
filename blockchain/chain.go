package blockchain

import (
	"container/list"
	"crypto/sha256"
)

type Blockchain struct {
	blocks *list.List
}

func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		blocks: list.New(),
	}

	return bc
}

func (bc *Blockchain) Append(b *Block) {
	bc.blocks.PushBack(b)
}

func ValidateBlockchain(bc *Blockchain) bool {
	// first block Previous is nil
	prev := [sha256.Size]byte{}
	first := true

	for e := bc.blocks.Front(); e != nil; e = e.Next() {
		b, ok := e.Value.(*Block)
		if !ok {
			return false
		}

		if !first && prev != b.Header.Prev {
			return false
		}

		c, err := b.Hash()
		if err != nil {
			return false
		}

		first = false
		prev = c
	}

	return true
}
