// Copyright 2017-present Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

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
