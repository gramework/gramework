package grypto

import (
	"github.com/gramework/gramework/grypto/providers/bcrypt"
	"github.com/gramework/gramework/grypto/providers/scrypt"
)

const defaultProvider = "scrypt"

var registry = map[string]PasswordCryptoProvider{
	"scrypt": scrypt.DefaultProvider,
	"2a":     bcrypt.DefaultProvider,
}
