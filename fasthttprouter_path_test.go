// Copyright 2013 Julien Schmidt. All rights reserved.
// Copyright (c) 2015-2016, 招牌疯子
// Copyright (c) 2017, Kirill Danshin
// Use of this source code is governed by a BSD-style license that can be found
// in the 3rd-Party License/fasthttprouter file.

package gramework

import "testing"

func TestCleanPath(t *testing.T) {
	if res := CleanPath(""); res != "/" {
		t.Errorf("expected: / actual: %s", res)
	}

	if res := CleanPath("/hello/../world"); res != "/world" {
		t.Errorf("expected: /world actual: %s", res)
	}
	if res := CleanPath("/hello/../../world"); res != "/world" {
		t.Errorf("expected: /world actual: %s", res)
	}
	if res := CleanPath("hello"); res != "/hello" {
		t.Errorf("expected: /hello actual: %s", res)
	}
	if res := CleanPath("hello/world"); res != "/hello/world" {
		t.Errorf("expected: /hello/world actual: %s", res)
	}
	if res := CleanPath("./hello/world"); res != "/hello/world" {
		t.Errorf("expected: /hello/world actual: %s", res)
	}
	if res := CleanPath("./hello/////world"); res != "/hello/world" {
		t.Errorf("expected: /hello/world actual: %s", res)
	}
	if res := CleanPath("./HeLLo/////world"); res != "/HeLLo/world" {
		t.Errorf("expected: /hello/world actual: %s", res)
	}
	if res := CleanPath("./hello/////world///abс//"); res != "/hello/world/abс/" {
		t.Errorf("expected: /hello/world/abc/ actual: %s", res)
	}
	if res := CleanPath("./hello/////world//../abс//"); res != "/hello/abс/" {
		t.Errorf("expected: /hello/abc/ actual: %s", res)
	}
}
