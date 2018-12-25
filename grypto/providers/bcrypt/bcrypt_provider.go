package bcrypt

import (
	"github.com/gramework/gramework"
	gobcrypt "golang.org/x/crypto/bcrypt"
)

var log = gramework.Logger.WithField("package", "gramework/grypto/providers/bcrypt")

const (
	DefaultCost = 13
)

var (
	DefaultProvider = New()

	providerName = []byte("2a")
)

// Provider handles algorythm parameters
type Provider struct {
	Cost uint8
}

// New returns new bcrypt provider
func New() *Provider {
	return &Provider{
		Cost: DefaultCost,
	}
}

// Hash returns scrypt hash of plaintext
func (p *Provider) Hash(plaintext []byte) []byte {
	pw, _ := gobcrypt.GenerateFromPassword(plaintext, int(p.Cost))
	return pw
}

// HashString returns scrypt hash of plaintext
func (p *Provider) HashString(plaintext string) []byte {
	return p.Hash([]byte(plaintext))
}

// NeedsRehash checks if provided hash needs rehash
func (p *Provider) NeedsRehash(hash []byte) bool {
	hashCost, err := gobcrypt.Cost(hash)
	return err != nil || hashCost != int(p.Cost)
}

// Valid checks if provided plaintext is valid for given hash
func (p *Provider) Valid(hash, plain []byte) bool {
	return gobcrypt.CompareHashAndPassword(hash, plain) == nil
}
