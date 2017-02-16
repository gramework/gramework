package apiClient

import (
	"errors"

	"github.com/valyala/bytebufferpool"
)

// ErrNoServerAvailable occurred when no server available in the pool
var ErrNoServerAvailable = errors.New("no server available")

var buffer bytebufferpool.Pool
