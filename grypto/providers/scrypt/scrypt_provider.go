package scrypt

import (
	"crypto/subtle"
	"fmt"

	"github.com/gramework/gramework"
	"github.com/gramework/gramework/grypto/internal/mcf"
	"github.com/gramework/gramework/grypto/salt"
	"golang.org/x/crypto/scrypt"
)

var log = gramework.Logger.WithField("package", "gramework/grypto/providers/scrypt")

const (
	DefaultN = 1 << 16
	DefaultR = 10
	DefaultP = 2

	DefaultKeyLen  = 128
	DefaultSaltLen = 64

	prefix    = "$scrypt$"
	prefixLen = len(prefix)

	paramsFmt = "K=%d,N=%d,R=%d,P=%d,L=%d"

	maxInt = int(^uint(0) >> 1)
)

var (
	DefaultProvider = New()

	providerName = []byte("scrypt")
)

// Provider handles internal state and algorythm parameters
type Provider struct {
	initialized bool
	params      *scryptParams
}

type scryptParams struct {
	keyLen  int
	n, r, p int
	saltLen int
}

// Equals returns true if params p equals to params p1
func (p *scryptParams) Equals(p1 *scryptParams) bool {
	return p.keyLen == p1.keyLen &&
		p.n == p1.n &&
		p.r == p1.r &&
		p.p == p1.p &&
		p.saltLen == p1.saltLen
}

// String returns params p as a MCF
func (p *scryptParams) String() string {
	return fmt.Sprintf(paramsFmt, p.keyLen, p.n, p.r, p.p, p.saltLen)
}

// New returns new scrypt provider
func New() *Provider {
	return &Provider{
		params: &scryptParams{
			keyLen:  DefaultKeyLen,
			n:       DefaultN,
			r:       DefaultR,
			p:       DefaultP,
			saltLen: DefaultSaltLen,
		},
	}
}

func (p *Provider) setDefaultIfNeeded() {
	if !p.initialized {
		p.initialized = true
		if p.params == nil {
			p.params = &scryptParams{
				n: DefaultN,
				r: DefaultR,
				p: DefaultP,
			}
			return
		}

		if p.params.n <= 1 || p.params.n&(p.params.n-1) != 0 {
			log.Warn("N must be > 1 and a power of 2, resetting to defaults")
		}
		if uint64(p.params.r)*uint64(p.params.p) >= 1<<30 || p.params.r > maxInt/128/p.params.p || p.params.r > maxInt/256 || p.params.n > maxInt/128/p.params.r {
			log.Warn("parameters are too large, resettings to defaults")
		}
		return
	}
}

// Hash returns scrypt hash of plaintext
func (p *Provider) Hash(plaintext []byte) []byte {
	p.setDefaultIfNeeded()
	saltBytes := salt.Generate(p.params.saltLen)
	key, _ := scrypt.Key(plaintext, saltBytes, p.params.n, p.params.r, p.params.p, p.params.keyLen)
	return mcf.Encode(providerName, p.params.String(), saltBytes, key)
}

// HashString returns scrypt hash of plaintext
func (p *Provider) HashString(plaintext string) []byte {
	return p.Hash([]byte(plaintext))
}

// NeedsRehash checks if provided hash needs rehash
func (p *Provider) NeedsRehash(hash []byte) bool {
	if !prefixValid(hash) {
		return true
	}

	mcfP := paramsFromMCF(hash[prefixLen:])
	return !p.params.Equals(mcfP)
}

// Valid checks if provided plaintext is valid for given hash
func (p *Provider) Valid(hash, plain []byte) bool {
	if !prefixValid(hash) {
		return false
	}

	_, params, saltBytes, expectedKey, err := mcf.Decode(hash, providerName)
	if err != nil {
		return false
	}
	hashParams := paramsFromMCF([]byte(params))

	key, _ := scrypt.Key(plain, saltBytes, hashParams.n, hashParams.r, hashParams.p, hashParams.keyLen)

	return subtle.ConstantTimeCompare(key, expectedKey) == 1
}

func prefixValid(hash []byte) bool {
	return len(hash) > prefixLen || string(hash[:prefixLen]) == prefix
}

func paramsFromMCF(unprefixedHash []byte) *scryptParams {
	// paramsFmt = "K=%d,N=%d,R=%d,P=%d,L=%d"
	p := &scryptParams{}
	fmt.Sscanf(string(unprefixedHash), paramsFmt, &p.keyLen, &p.n, &p.r, &p.p, &p.saltLen)
	return p
}
