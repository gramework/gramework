package test

import (
	"testing"

	"github.com/apex/log"
	"github.com/gramework/gramework"
)

func TestLogErrorfShouldNotPanic(t *testing.T) {
	defer func() {
		e := recover()
		if e != nil {
			t.Logf("panic handled while testing log.errorf: %+#v", e)
			t.FailNow()
		}
	}()
	gramework.Errorf("test: %v", []string{"test"})
}

func TestFastHTTPLoggerAdapter(t *testing.T) {
	var apexLogger log.Interface = gramework.Logger
	logger := gramework.NewFastHTTPLoggerAdapter(&apexLogger)
	logger.Printf("printed")
}
