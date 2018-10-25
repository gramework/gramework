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

	"github.com/stretchr/testify/assert"
)

func TestApp_SetName(t *testing.T) {
	setName := func(n string) {
		app := New(OptAppName(DefaultAppName))
		if assert.Equal(t, DefaultAppName, app.name) {
			app.SetName(n)
		}
		if len(n) == 0 {
			n = DefaultAppName
		}
		assert.Equal(t, n, app.name)
		assert.Equal(t, n, app.serverBase.Name)
	}
	t.Run("CustomName", func(t *testing.T) {
		setName("test_app")
	})
	t.Run("EmptyName", func(t *testing.T) {
		setName("")
	})
}
