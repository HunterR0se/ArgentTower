package game

import (
	"crypto/rand"
	"math/big"
)

// randFloat returns a random float64 in [0.0, 1.0)
func randFloat() float64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(1<<53))
	return float64(n.Int64()) / (1 << 53)
}