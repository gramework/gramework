package gramework

import "errors"

var (
	// ErrTLSNoEmails occurs when no emails provided but user tries to use AutoTLS features
	ErrTLSNoEmails = errors.New("auto tls: no emails provided")

	// ErrArgNotFound used when no route argument is found
	ErrArgNotFound = errors.New("undefined argument")
)
