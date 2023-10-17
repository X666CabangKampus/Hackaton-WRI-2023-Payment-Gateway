package signature

import (
	"crypto/rand"
	"crypto/rsa"
	"math/big"
)

const Bits int = 2048

func GenPairKey() (*big.Int, *big.Int) {
	key, _ := rsa.GenerateKey(rand.Reader, Bits)
	return key.N, key.D
}
