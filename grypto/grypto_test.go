package grypto

import (
	"testing"

	"github.com/gramework/gramework/grypto/providers/bcrypt"
	"github.com/gramework/gramework/grypto/providers/scrypt"
	"github.com/gramework/utils/grand"
)

func TestProviderMatchingAndSanity(t *testing.T) {
	p1 := scrypt.New()
	p2 := bcrypt.New()

	pw := make([]byte, 12)
	grand.Read(pw)

	phash1 := p1.Hash(pw)
	phash2 := p2.Hash(pw)

	if !p1.Valid(phash1, pw) {
		t.Error("p1 validation failed with original password")
		return
	}

	if !p2.Valid(phash2, pw) {
		t.Error("p2 validation failed with original password")
		return
	}

	if p1.Valid(phash2, pw) {
		t.Error("p1 validation succeed with original password but phash2")
		return
	}

	if p2.Valid(phash1, pw) {
		t.Error("p2 validation succeed with original password but phash1")
		return
	}

	if !PasswordValid(phash1, pw) {
		t.Error("phash1 matching validation failed with original password")
		return
	}

	if !PasswordValid(phash2, pw) {
		t.Error("phash2 matching validation failed with original password")
		return
	}

	var pwc = make([]byte, 12)
	copy(pwc, pw)

	pwc[2] ^= pwc[0]
	pwc[3] ^= pwc[2]

	if p1.Valid(phash1, pwc) {
		t.Error("p1 validation succeed with modified password")
		return
	}

	if p2.Valid(phash2, pwc) {
		t.Error("p2 validation succeed with modified password")
		return
	}

	if p1.Valid(phash2, pwc) {
		t.Error("p1 validation succeed with modified password but phash2")
		return
	}

	if p2.Valid(phash1, pwc) {
		t.Error("p2 validation succeed with modified password but phash1")
		return
	}
}
