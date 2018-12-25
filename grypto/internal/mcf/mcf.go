package mcf

import (
	"bytes"
	"crypto/subtle"
	"encoding/base64"
	"errors"
)

const Divider = '$'

var (
	ErrorInvalidDecodeInput = errors.New("gramework/grypto: invalid mcf input")

	splitDivider = []byte{Divider}
)

// Encode encodes given providerName, params, salt and key into Modular Crypt Format.
func Encode(providerName []byte, params string, salt []byte, key []byte) (res []byte) {
	salt64 := encodeBase64(salt)
	key64 := encodeBase64(key)
	// final res len = len(providerName) + len(params) + len(salt) + len(key) + 4, because we need 4 dividers
	res = append(res, Divider)
	res = append(res, providerName...)
	res = append(res, Divider)
	res = append(res, params...)
	res = append(res, Divider)
	res = append(res, salt64...)
	res = append(res, Divider)
	res = append(res, key64...)
	return res
}

// encodeBase64 encodes the input bytes into standard base64.
func encodeBase64(in []byte) (out []byte) {
	enc := base64.StdEncoding
	out = make([]byte, enc.EncodedLen(len(in)))
	enc.Encode(out, in)
	return out
}

func decodeBase64(in []byte) (out []byte, err error) {
	enc := base64.StdEncoding
	out = make([]byte, enc.DecodedLen(len(in)))
	n, err := enc.Decode(out, in)
	return out[:n], err
}

// Decode given MCF if it contains information about expected pw type
func Decode(mcf []byte, expectedPWType []byte) (providerName []byte, params string, salt []byte, key []byte, err error) {
	if len(mcf) <= 1 || mcf[0] != Divider {
		err = ErrorInvalidDecodeInput
		return
	}

	parts := bytes.Split(mcf[1:], splitDivider)

	if len(parts) != 4 {
		err = ErrorInvalidDecodeInput
		return
	}

	if subtle.ConstantTimeCompare(parts[0], expectedPWType) == 0 {
		err = ErrorInvalidDecodeInput
		return
	}

	params = string(parts[1])

	saltDecoded, err := decodeBase64(parts[2])
	if err != nil {
		err = ErrorInvalidDecodeInput
		return
	}

	salt = saltDecoded

	keyDecoded, err := decodeBase64(parts[3])
	if err != nil {
		err = ErrorInvalidDecodeInput
		return
	}
	key = keyDecoded
	return
}
