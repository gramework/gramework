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
	"os"
	"path/filepath"
	"testing"

	"github.com/valyala/fasthttp"
)

func testProvideTempFile(t *testing.T, action func(file, dir string)) {
	tmpDir := t.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, "servedirtmp_")
	if err != nil {
		t.Error("cannot create temporary file:", err)
	}

	// write some text into file
	text := []byte("test_file_data")
	if _, err = tmpFile.Write(text); err != nil {
		t.Error("failed to write to temporary file:", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Error(err)
	}

	action(filepath.Base(tmpFile.Name()), tmpDir)
}

func TestAppServeDir(t *testing.T) {
	var funcSD func(*Context)

	app := New()

	ctx := Context{RequestCtx: &fasthttp.RequestCtx{}}

	testProvideTempFile(t, func(file, dir string) {
		ctx.Request.SetRequestURI(file)
		funcSD = app.ServeDir(dir)
		funcSD(&ctx)
	})

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Errorf("failed on file read by directory server %d", ctx.Response.StatusCode())
	}
}

func TestAppServeDirNoCache(t *testing.T) {
	var funcSD func(*Context)
	const (
		hdrCacheControl = "no-cache, no-store, must-revalidate"
		hdrPragma       = "no-cache"
		hdrExpires      = "0"
	)

	app := New()

	ctx := Context{RequestCtx: &fasthttp.RequestCtx{}}

	testProvideTempFile(t, func(file, dir string) {
		ctx.Request.SetRequestURI(file)
		funcSD = app.ServeDirNoCache(dir)
		funcSD(&ctx)
	})

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Errorf("failed on directory serve (no cache) %d", ctx.Response.StatusCode())
	}

	hdrList := []struct {
		name string
		exp  string
	}{
		{"Cache-Control", hdrCacheControl},
		{"Pragma", hdrPragma},
		{"Expires", hdrExpires},
	}
	for _, hdr := range hdrList {
		if act := string(ctx.Response.Header.Peek(hdr.name)); act != hdr.exp {
			t.Errorf("wrong header '%s' check: %s, but %s expected", hdr.name, act, hdr.exp)
		}
	}
}
