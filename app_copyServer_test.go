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
	"reflect"
	"testing"
)

var fieldList = []string{
	"Handler",
	"Name",
	"Concurrency",
	"DisableKeepalive",
	"ReadBufferSize",
	"WriteBufferSize",
	"ReadTimeout",
	"WriteTimeout",
	"MaxConnsPerIP",
	"MaxRequestsPerConn",
	"MaxKeepaliveDuration",
	"MaxRequestBodySize",
	"ReduceMemoryUsage",
	"GetOnly",
	"LogAllErrors",
	"DisableHeaderNamesNormalizing",
	"Logger",
}

func compareField(t *testing.T, act, exp interface{}, field string) {
	actVal := reflect.Indirect(reflect.ValueOf(act)).FieldByName(field)
	expVal := reflect.Indirect(reflect.ValueOf(act)).FieldByName(field)
	if actVal != expVal {
		t.Errorf("field %s of app copy has value %v, but expected %v", field, actVal, expVal)
	}
}

func TestAppCopyServer(t *testing.T) {
	app := New()

	serverCopy := app.copyServer()

	for _, field := range fieldList {
		compareField(t, serverCopy, app.serverBase, field)
	}
}
