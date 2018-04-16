// Copyright 2017 Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

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
	return err != nil || hashCost != cost
}

// Salt128 generates 128 bits of random data.
func Salt128() []byte {
	x := make([]byte, 16)
	rand.Read(x)
	return x
}

// PasswordValid checks if provided hash
func PasswordValid(hash, password []byte) bool {
	return bcrypt.CompareHashAndPassword(hash, password) == nil
}
