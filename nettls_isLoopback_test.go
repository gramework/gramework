package gramework

import (
	"testing"
)

func TestIsLoopback(t *testing.T) {
	for k, v := range map[string]bool{
		"localhost":   true,
		"127.0.0.1":   true,
		"127.0.2.1":   true,
		"::1":         true,
		"example.com": false,
	} {
		if res := IsLoopback(k); res != v {
			t.Logf("IsLoopback test failed while testing [ %v ] case, expected: %v, got: %v", k, v, res)
			t.FailNow()
		}
	}
}
