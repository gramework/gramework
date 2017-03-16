package gramework

import (
	"testing"
)

func TestIsLoopback(t *testing.T) {
	for k, v := range map[string]bool{
		"localhost":         true,
		"127.0.0.1":         true,
		"127.0.2.1":         true,
		"::1":               true,
		"example.com":       false,
		"google.com":        false,
		"127.0.0.2.2.2.2.2": false,
		"3423423423423423":  false,
		"localhost:1234":    true,
		"localhost:":        true,
		"127.0.0.1:443":     true,
		"127.0.1.5":         true,
		"10.0.0.5":          false,
		"12.7.0.1":          false,
		"[::1]":             true,
		"[::1]:1234":        true,
		"::":                false,
		"[::]":              false,
		"local":             false,
		"У лупбека лупбек чудестный, лупбечна цепь на лупе том, лупом и беком лупбек лупбечный, все ходит сам в себя кругом": false, // UTF-8 test
	} {
		if res := IsLoopback(k); res != v {
			t.Logf("IsLoopback test failed while testing %q case, expected: %v, got: %v", k, v, res)
			t.FailNow()
		}
	}
}

func BenchmarkIsLoopbackBestCase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsLoopback(localhost)
	}
}

func BenchmarkIsLoopbackBestCase2(b *testing.B) {
	const addr = "localhost:443"
	for i := 0; i < b.N; i++ {
		IsLoopback(addr)
	}
}

func BenchmarkIsLoopbackMidCase(b *testing.B) {
	const addr = "example.com"
	for i := 0; i < b.N; i++ {
		IsLoopback(addr)
	}
}

func BenchmarkIsLoopbackWorstCase(b *testing.B) {
	const addr = "У лупбека лупбек чудестный, лупбечна цепь на лупе том, лупом и беком лупбек лупбечный, все ходит сам в себя кругом"
	for i := 0; i < b.N; i++ {
		IsLoopback(addr)
	}
}
