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
	"time"

	"github.com/valyala/fasthttp"
)

// ServeDir from a given path
func (app *App) ServeDir(path string) func(*Context) {
	return app.ServeDirCustom(path, 0, true, false, []string{"index.html", "index.htm"})
}

// ServeDirCustom gives you ability to serve a dir with custom settings
func (app *App) ServeDirCustom(path string, stripSlashes int, compress bool, generateIndexPages bool, indexNames []string) func(*Context) {
	if indexNames == nil {
		indexNames = []string{}
	}
	fs := &fasthttp.FS{
		Root:                 path,
		IndexNames:           indexNames,
		GenerateIndexPages:   generateIndexPages,
		Compress:             compress,
		CacheDuration:        5 * time.Minute,
		CompressedFileSuffix: ".gz",
	}

	if stripSlashes > 0 {
		fs.PathRewrite = fasthttp.NewPathSlashesStripper(stripSlashes)
	}

	h := fs.NewRequestHandler()
	return func(ctx *Context) {
		h(ctx.RequestCtx)
	}
}

// ServeDirNoCache gives you ability to serve a dir without caching
func (app *App) ServeDirNoCache(path string) func(*Context) {
	return app.ServeDirNoCacheCustom(path, 0, true, false, nil)
}

// ServeDirNoCacheCustom gives you ability to serve a dir with custom settings without caching
func (app *App) ServeDirNoCacheCustom(path string, stripSlashes int, compress bool, generateIndexPages bool, indexNames []string) func(*Context) {
	if indexNames == nil {
		indexNames = []string{}
	}
	fs := &fasthttp.FS{
		Root:                 path,
		IndexNames:           indexNames,
		GenerateIndexPages:   generateIndexPages,
		Compress:             compress,
		CacheDuration:        time.Millisecond,
		CompressedFileSuffix: ".gz",
	}

	if stripSlashes > 0 {
		fs.PathRewrite = fasthttp.NewPathSlashesStripper(stripSlashes)
	}

	h := fs.NewRequestHandler()
	pragmaH := "Pragma"
	pragmaV := "no-cache"
	expiresH := "Expires"
	expiresV := "0"
	ccH := "Cache-Control"
	ccV := "no-cache, no-store, must-revalidate"
	return func(ctx *Context) {
		ctx.Response.Header.Add(pragmaH, pragmaV)
		ctx.Response.Header.Add(expiresH, expiresV)
		ctx.Response.Header.Add(ccH, ccV)
		h(ctx.RequestCtx)
	}
}
