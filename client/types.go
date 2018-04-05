package client

import (
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

type (
	// Instance handles internal client representation and settings
	Instance struct {
		conf      *Config
		clients   map[string]*fasthttp.HostClient
		clientsMu *sync.RWMutex
		balancer  *rangeBalancer
	}

	// Config handles APIClient parameters
	Config struct {
		Addresses       []string
		WatcherTickTime time.Duration
	}

	rangeBalancer struct {
		total *int64
		curr  *int64
	}

	requestInfo struct {
		HostClient *fasthttp.HostClient
		Addr       string
	}
)
