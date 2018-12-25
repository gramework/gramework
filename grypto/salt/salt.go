package salt

import "crypto/rand"

func Generate(bytes int) []byte {
	x := make([]byte, bytes)
	rand.Read(x)
	return x
}

func Gen128() []byte {
	return Generate(16)
}
