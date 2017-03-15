package gramework

import "testing"

func TestNewShouldNeverReturnNil(t *testing.T) {
	app := New()
	if app == nil {
		t.Log("App is nil!")
		t.FailNow()
		return
	}
	if app.defaultRouter == nil {
		t.Log("App Router is nil!")
		t.FailNow()
		return
	}
	if app.Logger == nil {
		t.Log("App Logger is nil!")
		t.FailNow()
		return
	}
	if app.domainListLock == nil {
		t.Log("App domain list lock is nil!")
		t.FailNow()
		return
	}
}
