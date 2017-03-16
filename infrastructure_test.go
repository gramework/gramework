package gramework

import (
	"testing"

	"github.com/gramework/gramework/infrastructure"
)

func TestServeInfrastructure(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Logf("Serve Infrastructure should not panic, but %v", e)
			t.FailNow()
		}
	}()
	app := New()
	i := infrastructure.New()
	app.ServeInfrastructure(i)
}
