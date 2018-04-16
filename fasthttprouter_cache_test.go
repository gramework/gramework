// Copyright 2017 Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package gramework

import (
	"testing"
)

func TestRouterCache(t *testing.T) {
	cache := &cache{
		v: make(map[string]*msc, zero),
	}

	if _, ok := cache.Get(Slash, GET); ok {
		t.Fatalf("Cache returned ok flag for key that not exists")
	}

	cache.Put(Slash, new(node), false, GET)

	if n, ok := cache.Get(Slash, GET); !ok || n == nil {
		t.Fatalf("Cache returned unexpected result: n=[%v], ok=[%v]", n, ok)
	}
}
