package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"math/bits"
)

const (
	ProofOfWorkDifficulty int = 9
	SignatureSize         int = 91
)

type Header struct {
	Prev  [sha256.Size]byte
	Nonce [sha256.Size]byte
	Root  [sha256.Size]byte
}

type Block struct {
	Header Header
	Owner  [91]byte
}

// DecodeBlock from binary format
func DecodeBlock(data []byte) (*Block, error) {
	buf := bytes.NewBuffer(data)
	b := &Block{}

	err := binary.Read(buf, binary.LittleEndian, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Encode block into binary format
func (b Block) Encode() ([]byte, error) {
	var buf bytes.Buffer

	err := binary.Write(&buf, binary.LittleEndian, b)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (b Block) Hash() ([sha256.Size]byte, error) {

	bin, err := b.Encode()
	if err != nil {
		return sha256.Sum256(nil), err
	}

	return sha256.Sum256(bin), nil
}

func Validate(b *Block) (bool, error) {
	h, err := b.Hash()
	if err != nil {
		return false, err
	}

	var c int
	for _, byt := range h {
		i := bits.LeadingZeros8(byt)
		c += i
		if i != 8 {
			break
		}
	}

	fmt.Printf("%x ---> c=%d\n", h, c)
	if c < ProofOfWorkDifficulty {
		return false, nil
	}

	return true, nil
}

func GenerateValideHash(b *Block) error {
	buf := make([]byte, 512)

	for {
		_, err := rand.Read(buf)
		if err != nil {
			return err
		}
		b.Header.Nonce = sha256.Sum256(buf)
		ok, err := Validate(b)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}
}

// Sign previous hash and insert public key
func (b *Block) Sign(privateKey *ecdsa.PrivateKey) error {
	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, b.Header.Prev[:])
	if err != nil {
		return fmt.Errorf("cannot sign previous block's hash: %s", err)
	}
	b.Header.Prev = [sha256.Size]byte(sig)

	ecdshPublicKey, err := privateKey.PublicKey.ECDH()
	if err != nil {
		return fmt.Errorf("cannot convert public key to ECDH: %s", err)
	}

	pubkey, err := x509.MarshalPKIXPublicKey(ecdshPublicKey)
	if err != nil {
		return fmt.Errorf("cannot encode public key to ASN1 DER: %s", err)
	}

	b.Owner = [91]byte(pubkey)
	return nil
}
