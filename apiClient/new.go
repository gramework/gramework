package apiClient

import (
	"sync"

	"github.com/valyala/fasthttp"
)

// New API client instance
func New(config Config) *Instance {
	client := &Instance{
		conf:      &config,
		clients:   make(map[string]*fasthttp.HostClient),
		clientsMu: &sync.RWMutex{},
		balancer:  newRangeBalancer(),
	}

	return client
}
