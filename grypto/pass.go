// Package grypto provides helpers for dealing with cryptography
package grypto

import (
	"crypto/rand"

	"golang.org/x/crypto/bcrypt"
)

const (
	cost = 10
)

// PasswordHashString returns hash of plain password in the given string
func PasswordHashString(plainPass string) []byte {
	return PasswordHash([]byte(plainPass))
}

// PasswordHash returns hash of plain password in the given byte slice
func PasswordHash(plainPass []byte) []byte {
	pw, _ := bcrypt.GenerateFromPassword(plainPass, cost)
	return pw
}

// PasswordNeedsRehash checks if the password should be rehashed as soon as possible
func PasswordNeedsRehash(hash []byte) bool {
	hashCost, err := bcrypt.Cost(hash)
	if err != nil || hashCost != cost {
		return true
	}
	return false
}

// Salt128 generates 128 bits of random data.
func Salt128() []byte {
	x := make([]byte, 16)
	rand.Read(x)
	return x
}

// PasswordValid checks if provided hash
func PasswordValid(hash, password []byte) bool {
	if err := bcrypt.CompareHashAndPassword(hash, password); err != nil {
		return false
	}
	return true
}
