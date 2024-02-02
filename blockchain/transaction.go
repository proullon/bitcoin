package blockchain

type Transaction struct {
	In  [SignatureSize]byte
	Out map[[SignatureSize]byte]float64
}
