// Copyright 2017-present Kirill Danshin and Gramework contributors
// Copyright 2019-present Highload LTD (UK CN: 11893420)
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

	"github.com/stretchr/testify/assert"
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

func TestNewWithName(t *testing.T) {
	newApp := func(n string) {
		app := New(OptAppName(n))
		assert.Equal(t, n, app.name)
		assert.Equal(t, n, app.serverBase.Name)
	}
	t.Run("DefaultName", func(t *testing.T) {
		newApp(DefaultAppName)
	})
	t.Run("CustomName", func(t *testing.T) {
		newApp("test_app")
	})
	t.Run("EmptyName", func(t *testing.T) {
		newApp("")
	})
}

func TestNewWithoutName(t *testing.T) {
	app := New()
	assert.Equal(t, DefaultAppName, app.name)
	assert.Equal(t, DefaultAppName, app.serverBase.Name)
}
