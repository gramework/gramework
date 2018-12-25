// Copyright 2017-present Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

// Package grypto provides helpers for dealing with cryptography.
package grypto

import (
	"crypto/subtle"

	"github.com/gramework/gramework/grypto/salt"
)

type PasswordCryptoProvider interface {
	HashString(plain string) []byte
	Hash(plain []byte) []byte
	NeedsRehash(hash []byte) bool
	Valid(hash, plain []byte) bool
}

func matchProvider(hash []byte) PasswordCryptoProvider {
	if len(hash) < 4 {
		return nil
	}
	for key, provider := range registry {
		if len(hash) > len(key)+2 && subtle.ConstantTimeCompare(hash[1:len(key)+1], []byte(key)) == 1 {
			return provider
		}
	}
	return nil
}

// PasswordHashString returns hash of plain password in the given string
func PasswordHashString(plainPass string) []byte {
	return registry[defaultProvider].HashString(plainPass)
}

// PasswordHash returns hash of plain password in the given byte slice
func PasswordHash(plainPass []byte) []byte {
	return registry[defaultProvider].Hash(plainPass)
}

// PasswordNeedsRehash checks if the password should be rehashed as soon as possible
func PasswordNeedsRehash(hash []byte) bool {
	return registry[defaultProvider].NeedsRehash(hash)
}

// Salt128 generates 128 bits of random data.
func Salt128() []byte {
	return salt.Gen128()
}

// PasswordValid checks if provided hash
func PasswordValid(hash, password []byte) bool {
	p := matchProvider(hash)
	if p == nil {
		return false
	}
	return p.Valid(hash, password)
}
