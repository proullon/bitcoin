package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"testing"
)

func _testHash() [sha256.Size]byte {
	data := []byte("Hello world!")
	return sha256.Sum256(data)
}

func TestBlockHash(t *testing.T) {
	var b Block
	b.Header.Prev = _testHash()
	b.Header.Nonce = _testHash()
	b.Header.Root = _testHash()

	h, err := b.Hash()
	if err != nil {
		t.Fatalf("error while encoding: %s", err)
	}

	t.Logf("Block hash: %d", h)
}

func TestValidity(t *testing.T) {
	b := &Block{}
	b.Header.Prev = _testHash()
	b.Header.Nonce = _testHash()
	b.Header.Root = _testHash()

	ok, err := Validate(b)
	if err != nil {
		t.Fatalf("unexpected error validating block: %s", err)
	}
	if ok {
		t.Fatalf("expected invalid proof of work")
	}

	err = GenerateValideHash(b)
	if err != nil {
		t.Fatalf("unexpected error generating proof of work")
	}

	ok, err = Validate(b)
	if err != nil {
		t.Fatalf("unexpected error validating block: %s", err)
	}
	if !ok {
		t.Fatalf("expected valid proof of work")
	}

}

func TestSignature(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("unexpected error generating private key: %s", err)
	}

	b := &Block{}
	prev := _testHash()
	b.Header.Nonce = _testHash()
	b.Header.Root = _testHash()

	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, prev[:])
	if err != nil {
		t.Fatalf("cannot sign previous block's hash: %s", err)
	}
	t.Logf("Signature: %x", sig)
	b.Header.Prev = [sha256.Size]byte(sig)

	ecdshPublicKey, err := privateKey.PublicKey.ECDH()
	if err != nil {
		t.Fatalf("cannot convert public key to ECDH: %s", err)
	}

	pubkey, err := x509.MarshalPKIXPublicKey(ecdshPublicKey)
	if err != nil {
		t.Fatalf("cannot encode public key to ASN1 DER: %s", err)
	}

	t.Logf("Pubkey: %x  (len=%d)", pubkey, len(pubkey))
	b.Owner = [91]byte(pubkey)
}
