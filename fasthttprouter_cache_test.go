package gramework

import "testing"

func TestRouterCache(t *testing.T) {
	cache := &cache{
		v: make(map[string]*msc, 0),
	}

	if _, ok := cache.Get("/", GET); ok {
		t.Fatalf("Cache returned ok flag for key that not exists")
	}

	cache.Put("/", &node{}, false, GET)

	if n, ok := cache.Get("/", GET); !ok || n == nil {
		t.Fatalf("Cache returned unexpected result: n=[%v], ok=[%v]", n, ok)
	}
}
