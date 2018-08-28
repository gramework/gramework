// Copyright 2017-present Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package gramework

import (
	"testing"
)

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
