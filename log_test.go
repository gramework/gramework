package gramework

import "testing"

func TestLogErrorfShouldNotPanic(t *testing.T) {
	defer func() {
		e := recover()
		if e != nil {
			t.Logf("panic handled while testing log.errorf: %+#v", e)
			t.FailNow()
		}
	}()
	Errorf("test: %v", []string{"test"})
}
