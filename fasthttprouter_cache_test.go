package gramework

import "testing"

func TestRouterCache(t *testing.T) {
	cache := &cache{
		v: make(map[string]*msc, zero),
	}

	if _, ok := cache.Get(Slash, GET); ok {
		t.Fatalf("Cache returned ok flag for key that not exists")
	}

	cache.Put(Slash, &node{}, false, GET)

	if n, ok := cache.Get(Slash, GET); !ok || n == nil {
		t.Fatalf("Cache returned unexpected result: n=[%v], ok=[%v]", n, ok)
	}
}
