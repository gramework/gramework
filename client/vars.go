package apiClient

import (
	"errors"

	"github.com/valyala/bytebufferpool"
)

var (
	// ErrNoServerAvailable occurred when no server available in the pool
	ErrNoServerAvailable = errors.New("no server available")

	buffer bytebufferpool.Pool
)
