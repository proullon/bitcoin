package blockchain

import (
	"crypto/sha256"
	"testing"
)

func TestValidate(t *testing.T) {

	bc := NewBlockchain()

	ok := ValidateBlockchain(bc)
	if !ok {
		t.Fatalf("expected empty blockchain to be valid")
	}

	b := &Block{
		Header: Header{
			Prev:  [sha256.Size]byte{},
			Nonce: [sha256.Size]byte{26, 40, 18, 41, 246, 84, 101, 88, 229, 243, 242, 66, 210, 31, 237, 17, 78, 78, 6, 29, 111, 30, 182, 10, 187, 101, 182, 53, 188, 21, 217, 201},
			Root:  [sha256.Size]byte{26, 40, 18, 41, 246, 84, 101, 88, 229, 243, 242, 66, 210, 31, 237, 17, 78, 78, 6, 29, 111, 30, 182, 10, 187, 101, 182, 53, 188, 21, 217, 201},
		},
	}
	bc.Append(b)

	ok = ValidateBlockchain(bc)
	if !ok {
		t.Fatalf("expected 1 block blockchain to be valid")
	}

	prev, err := b.Hash()
	if err != nil {
		t.Fatalf("cannot hash block: %s", err)
	}

	b = &Block{
		Header: Header{
			Prev:  prev,
			Nonce: [sha256.Size]byte{26, 40, 18, 41, 246, 84, 101, 88, 229, 243, 242, 66, 210, 31, 237, 17, 78, 78, 6, 29, 111, 30, 182, 10, 187, 101, 182, 53, 188, 21, 217, 201},
			Root:  [sha256.Size]byte{26, 40, 18, 41, 246, 84, 101, 88, 229, 243, 242, 66, 210, 31, 237, 17, 78, 78, 6, 29, 111, 30, 182, 10, 187, 101, 182, 53, 188, 21, 217, 201},
		},
	}
	bc.Append(b)
	ok = ValidateBlockchain(bc)
	if !ok {
		t.Fatalf("expected 2 block blockchain to be valid")
	}

	bc.Append(b)
	ok = ValidateBlockchain(bc)
	if ok {
		t.Fatalf("expected invalid previous block")
	}
}
