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
		"У лупбека лупбек чудестный, лупбечна цепь на лупе том, лупом и беком лупбек лупбечный, все ходит сам в себя кругом": false,
	} {
		if res := IsLoopback(k); res != v {
			t.Logf("IsLoopback test failed while testing [ %v ] case, expected: %v, got: %v", k, v, res)
			t.FailNow()
		}
	}
}
