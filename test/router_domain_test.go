package test

import (
	"testing"

	"github.com/gramework/gramework"
)

func TestDomainShouldNeverReturnNil(t *testing.T) {
	app := gramework.New()
	if app.Domain("test") == nil {
		t.FailNow()
	}
}
