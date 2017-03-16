package gramework

import "testing"

func TestDomainShouldNeverReturnNil(t *testing.T) {
	app := New()
	if app.Domain("test") == nil {
		t.FailNow()
	}
}
