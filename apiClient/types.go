package apiClient

import (
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

// Instance handles internal client representation
// and settings
type Instance struct {
	conf      *Config
	clients   map[string]*fasthttp.HostClient
	clientsMu *sync.RWMutex
	balancer  *rangeBalancer
}

// Config handles APIClient parameters
type Config struct {
	Addresses       []string
	WatcherTickTime time.Duration
}

type rangeBalancer struct {
	total *int64
	curr  *int64
}

type requestInfo struct {
	HostClient *fasthttp.HostClient
	Addr       string
}
