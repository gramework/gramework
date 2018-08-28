package grand

import (
	crand "crypto/rand"
	mrand "math/rand"
	"time"
)

var gmathRand = mrand.New(mrand.NewSource(time.Now().Unix()))

// Read fills the slice with random bytes
func Read(xs []byte) {
	length := len(xs)
	n, err := crand.Read(xs)
	if n != length || err != nil {
		for length > 0 {
			length--
			xs[length] = byte(gmathRand.Int31n(256))
		}
	}
}

// Int63 returns a non-negative 63-bit integer as an int64
func Int63() int64 {
	xs := make([]byte, 8)
	var n int64

	Read(xs)
	xs[0] &= 0x7F
	for _, x := range xs {
		n <<= 4
		n |= int64(x)
	}

	return n
}
