package test

import (
	"testing"

	"github.com/gramework/utils/grand"
	"github.com/toby3d/gramework/grypto"
	"golang.org/x/crypto/bcrypt"
)

const (
	cost = 10
)

// TestSalt128 makes sure the result is always 32 hex characters
func TestSalt128(t *testing.T) {
	var salt []byte
	for i := 0; i < 1024; i++ {
		salt = grypto.Salt128()
		if len(salt) != 16 {
			t.Errorf("Salt128 fail to generate 128 bit salt: %02x", salt)
			t.FailNow()
		}
	}
}

// TestPasswordSanity makes sure PasswordHash and ComparePassword work well together
func TestPasswordSanity(t *testing.T) {
	pw := make([]byte, 12)
	pw2 := make([]byte, 12)
	var hash, hash2 []byte

	for i := 0; i < 32; i++ {
		grand.Read(pw)
		grand.Read(pw2)
		hash = grypto.PasswordHash(pw)
		hash2 = grypto.PasswordHash(pw2)

		if !grypto.PasswordValid(hash, pw) {
			t.Errorf("PasswordValid should return true for the pair: %s and %s", hash, pw)
			t.FailNow()
		}
		if grypto.PasswordValid(hash2, pw) {
			t.Error("PasswordValid is giving false positive")
			t.FailNow()
		}
	}
}

// TestPasswordStringSanity makes sure PasswordHashString and ComparePassword work well together
func TestPasswordStringSanity(t *testing.T) {
	pw := make([]byte, 12)
	pw2 := make([]byte, 12)
	var hash, hash2 []byte

	for i := 0; i < 32; i++ {
		grand.Read(pw)
		grand.Read(pw2)
		hash = grypto.PasswordHashString(string(pw))
		hash2 = grypto.PasswordHashString(string(pw2))

		if !grypto.PasswordValid(hash, pw) {
			t.Errorf("PasswordValid should return true for the pair: %s and %s", hash, pw)
			t.FailNow()
		}
		if grypto.PasswordValid(hash2, pw) {
			t.Error("PasswordValid is giving false positive")
			t.FailNow()
		}
	}
}

// TestPasswordNeedsRehash makes sure TestPasswordNeedsRehash works well
func TestPasswordNeedsRehash(t *testing.T) {
	pw := make([]byte, 12)
	for i := 0; i < 32; i++ {
		grand.Read(pw)
		hash, _ := bcrypt.GenerateFromPassword(pw, cost-1)
		if !grypto.PasswordNeedsRehash(hash) {
			t.Errorf("PasswordNeedsRehash returned false, expected true")
			t.FailNow()
		}
	}
	for i := 0; i < 32; i++ {
		grand.Read(pw)
		hash, _ := bcrypt.GenerateFromPassword(pw, cost)
		if grypto.PasswordNeedsRehash(hash) {
			t.Errorf("PasswordNeedsRehash returned true, expected false")
		}
	}
}

func BenchmarkPassHashAndValidation(b *testing.B) {
	pws := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		pws[i] = make([]byte, 12)
		grand.Read(pws[i])
	}

	b.SetBytes(0)
	b.ResetTimer()

	for i := 0; i < b.N; i += 2 {
		hash := grypto.PasswordHash(pws[i])
		if !grypto.PasswordValid(hash, pws[i]) {
			b.Errorf("PasswordValid should return true for the pair: %s and %s", hash, pws[i])
			b.FailNow()
		}
	}
}

func BenchmarkPassHash(b *testing.B) {
	pws := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		pws[i] = make([]byte, 12)
		grand.Read(pws[i])
	}

	b.SetBytes(0)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hash := grypto.PasswordHash(pws[i])
		_ = hash
	}
}
