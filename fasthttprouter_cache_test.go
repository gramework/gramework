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
)

func TestRouterCache(t *testing.T) {
	cache := &cache{
		v: make(map[string]*msc, zero),
	}

	if _, ok := cache.Get(Slash, GET); ok {
		t.Fatalf("Cache returned ok flag for key that not exists")
	}

	cache.Put(Slash, new(node), false, GET)

	if n, ok := cache.Get(Slash, GET); !ok || n == nil {
		t.Fatalf("Cache returned unexpected result: n=[%v], ok=[%v]", n, ok)
	}
}

// func testBuildReqRes(method, uri string) (*fasthttp.Request, *fasthttp.Response) {
// 	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
// 	req.Header.SetMethod(method)
// 	req.SetRequestURI(uri)
// 	return req, res
// }

// func TestAppServe(t *testing.T) {
// 	const uri = "http://test.request"

// 	testCases := []func(*App) (func(string, interface{}) *App, string){
// 		// check GET request
// 		func(app *App) (func(string, interface{}) *App, string) {
// 			return app.GET, GET
// 		},
// 		// check POST request
// 		func(app *App) (func(string, interface{}) *App, string) {
// 			return app.POST, POST
// 		},
// 		// check PUT request
// 		func(app *App) (func(string, interface{}) *App, string) {
// 			return app.PUT, PUT
// 		},
// 		// check PATCH request
// 		func(app *App) (func(string, interface{}) *App, string) {
// 			return app.PATCH, PATCH
// 		},
// 		// check DELETE request
// 		func(app *App) (func(string, interface{}) *App, string) {
// 			return app.DELETE, DELETE
// 		},
// 		// check HEAD request
// 		func(app *App) (func(string, interface{}) *App, string) {
// 			return app.HEAD, HEAD
// 		},
// 		// check OPTIONS request
// 		func(app *App) (func(string, interface{}) *App, string) {
// 			return app.OPTIONS, OPTIONS
// 		},
// 	}

// 	for _, test := range testCases {
// 		var handleOK bool

// 		app := New()
// 		ln := fasthttputil.NewInmemoryListener()
// 		c := &fasthttp.Client{
// 			Dial: func(addr string) (net.Conn, error) {
// 				return ln.Dial()
// 			},
// 		}

// 		reg, method := test(app)

// 		go func() {
// 			_ = app.Serve(ln)
// 		}()

// 		reg("/", func() {
// 			handleOK = true
// 		})
// 		err := c.Do(testBuildReqRes(method, uri))
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		ln.Close()

// 		if !handleOK {
// 			t.Errorf("%s request was not served correctly", method)
// 		}
// 	}
// }
