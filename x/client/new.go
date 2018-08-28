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

	"github.com/valyala/fasthttp"
)

// New API client instance
func New(config Config) *Instance {
	client := &Instance{
		conf:      &config,
		clients:   make(map[string]*fasthttp.HostClient),
		clientsMu: new(sync.RWMutex),
		balancer:  newRangeBalancer(),
	}

	return client
}
