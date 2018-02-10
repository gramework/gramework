package gramework

import (
	"github.com/apex/log"
	"testing"
)

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

func TestFastHTTPLoggerAdapter(t *testing.T) {
	var apexLogger log.Interface
	apexLogger = Logger
	logger := NewFastHTTPLoggerAdapter(&apexLogger)
	logger.Printf("printed")
}
