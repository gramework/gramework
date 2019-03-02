package salt

import (
	"crypto/rand"
	"errors"
)

var nonNilErr = errors.New("<placeholder>")

func Generate(bytes int) []byte {
	x := make([]byte, bytes)

	err := nonNilErr
	for err != nil {
		_, err = rand.Read(x)
	}

	return x
}

func Gen128() []byte {
	return Generate(16)
}
