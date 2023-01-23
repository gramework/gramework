// Copyright 2013 Julien Schmidt. All rights reserved.
// Copyright (c) 2015-2016, 招牌疯子
// Copyright (c) 2017, Kirill Danshin
// Use of this source code is governed by a BSD-style license that can be found
// in the 3rd-Party License/fasthttprouter file.

package gramework

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
)

type (
	// handlerStruct struct {
	// 	handeled *bool
	// }

	// mockFileSystem struct {
	// 	opened bool
	// }

	readWriter struct {
		net.Conn
		r bytes.Buffer
		w bytes.Buffer
	}
)

var zeroTCPAddr = &net.TCPAddr{
	IP: net.IPv4zero,
}

func TestRouter(t *testing.T) {
	router := New()

	routed := false
	router.Handle("GET", "/user/:name", func(ctx *Context) {
		routed = true
		want := map[string]string{"name": "gopher"}

		if ctx.UserValue("name") != want["name"] {
			t.Fatalf("wrong wildcard values: want %v, got %v", want["name"], ctx.UserValue("name"))
		}
		ctx.Success("foo/bar", []byte("success"))
	})

	s := &fasthttp.Server{
		Handler: router.handler(),
	}

	rw := new(readWriter)
	rw.r.WriteString("GET /user/gopher?baz HTTP/1.1\r\n\r\n")

	ch := make(chan error)
	go func() {
		ch <- s.ServeConn(rw)
	}()

	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("return error %s", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout")
	}

	if !routed {
		t.Fatal("routing failed")
	}
}

// func (h handlerStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	*h.handeled = true
// }

func TestRouterAPI(t *testing.T) {
	var get, head, options, post, put, patch, deleted bool

	app := New()
	app.GET("/GET", func(ctx *Context) {
		get = true
	})
	app.HEAD("/GET", func(ctx *Context) {
		head = true
	})
	app.OPTIONS("/GET", func(ctx *Context) {
		options = true
	})
	app.POST("/POST", func(ctx *Context) {
		post = true
	})
	app.PUT("/PUT", func(ctx *Context) {
		put = true
	})
	app.PATCH("/PATCH", func(ctx *Context) {
		patch = true
	})
	app.DELETE("/DELETE", func(ctx *Context) {
		deleted = true
	})

	s := &fasthttp.Server{
		Handler: app.handler(),
	}

	rw := new(readWriter)
	ch := make(chan error)

	rw.r.WriteString("GET /GET HTTP/1.1\r\n\r\n")
	go func() {
		ch <- s.ServeConn(rw)
	}()
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("return error %s", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout")
	}
	if !get {
		t.Error("routing GET failed")
	}

	rw.r.WriteString("HEAD /GET HTTP/1.1\r\n\r\n")
	go func() {
		ch <- s.ServeConn(rw)
	}()
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("return error %s", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout")
	}
	if !head {
		t.Error("routing HEAD failed")
	}

	rw.r.WriteString("OPTIONS /GET HTTP/1.1\r\n\r\n")
	go func() {
		ch <- s.ServeConn(rw)
	}()
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("return error %s", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout")
	}
	if !options {
		t.Error("routing OPTIONS failed")
	}

	rw.r.WriteString("POST /POST HTTP/1.1\r\n\r\n")
	go func() {
		ch <- s.ServeConn(rw)
	}()
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("return error %s", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout")
	}
	if !post {
		t.Error("routing POST failed")
	}

	rw.r.WriteString("PUT /PUT HTTP/1.1\r\n\r\n")
	go func() {
		ch <- s.ServeConn(rw)
	}()
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("return error %s", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout")
	}
	if !put {
		t.Error("routing PUT failed")
	}

	rw.r.WriteString("PATCH /PATCH HTTP/1.1\r\n\r\n")
	go func() {
		ch <- s.ServeConn(rw)
	}()
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("return error %s", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout")
	}
	if !patch {
		t.Error("routing PATCH failed")
	}

	rw.r.WriteString("DELETE /DELETE HTTP/1.1\r\n\r\n")
	go func() {
		ch <- s.ServeConn(rw)
	}()
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("return error %s", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout")
	}
	if !deleted {
		t.Error("routing DELETE failed")
	}
}

func TestRouterWildAnyCache(t *testing.T) {
	var get bool

	app := New()
	app.GET("/*any", func(ctx *Context) {
		get = true
	})
	s := &fasthttp.Server{
		Handler: app.handler(),
	}

	rw := new(readWriter)
	ch := make(chan error)

	boilCache := 64 // It works on values greater than the cache threshold (32).
	for i := 0; i < boilCache; i++ {
		app.defaultRouter.Allowed("/GET", "OPTIONS")
	}

	rw.r.WriteString("GET /GET HTTP/1.1\r\n\r\n")
	go func() {
		ch <- s.ServeConn(rw)
	}()
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("return error %s", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout")
	}
	if !get {
		t.Error("routing GET failed")
	}
}

func TestRouterWildAnyWithArgsCache(t *testing.T) {
	var get bool
	var getArgs map[string][]string
	var green bool
	var greenValue string

	app := New()
	app.GET("/GET/*any", func(ctx *Context) {
		get = true
		getArgs = ctx.GETParams()
	})
	app.GET("/GREEN/:CORN", func(ctx *Context) {
		green = true
		greenValue = ctx.UserValue("CORN").(string)
	})

	s := &fasthttp.Server{
		Handler: app.handler(),
	}

	rw := new(readWriter)
	ch := make(chan error)

	boilCache := 64 // It works on values greater than the cache threshold (32).
	for i := 0; i < boilCache; i++ {
		app.defaultRouter.Allowed("/GET/ANY?foo=bar&fizz=buzz&fish", "OPTIONS")
		app.defaultRouter.Allowed("/GET/ANY", "OPTIONS")
		app.defaultRouter.Allowed("/GET", "OPTIONS")
		app.defaultRouter.Allowed("/GREEN/CORN", "OPTIONS")
		app.defaultRouter.Allowed("/GREEN", "OPTIONS")
	}

	rw.r.WriteString("GET /GET/ANY?foo=bar&fizz=buzz&fish HTTP/1.1\r\n\r\n")
	go func() {
		ch <- s.ServeConn(rw)
	}()
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("return error %s", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout")
	}
	if !get {
		t.Error("routing GET failed")
	}
	if getArgs == nil {
		t.Error("get args should not be empty")
	}
	if foo, ok := getArgs["foo"]; ok {
		if len(foo) != 1 || foo[0] != "bar" {
			t.Error("the foo arg lost its bar value")
		}
	} else {
		t.Error("the foo arg must exist")
	}
	if fizz, ok := getArgs["fizz"]; ok {
		if len(fizz) != 1 || fizz[0] != "buzz" {
			t.Error("the fizz arg lost its buzz value")
		}
	} else {
		t.Error("the fizz arg must exist")
	}
	if fish, ok := getArgs["fish"]; ok {
		if len(fish) > 0 && len(strings.Join(fish, "")) > 0 {
			t.Error("how much is the fish?", fish)
		}
	} else {
		t.Error("the fish arg must exist")
	}

	rw.r.WriteString("GET /GREEN/CORN HTTP/1.1\r\n\r\n")
	go func() {
		ch <- s.ServeConn(rw)
	}()
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("return error %s", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout")
	}
	if !green {
		t.Error("routing GREEN failed")
	}
	if len(greenValue) == 0 {
		t.Error("the green value must exist")
	}
	if greenValue != "CORN" {
		t.Errorf("the green value must be equal to 'CORN', but have '%s'", greenValue)
	}
}

func TestRouterRoot(t *testing.T) {
	router := New()
	recv := catchPanic(func() {
		router.GET("noSlashRoot", nil)
	})
	if recv == nil {
		t.Fatal("registering path not beginning with '/' did not panic")
	}
}

// func TestRouterOPTIONS(t *testing.T) {
// 	// TODO: because fasthttp is not support OPTIONS method now,
// 	// these test cases will be used in the future.
// 	handlerFunc := func(_ *fasthttp.RequestCtx) {}

// 	router := New()
// 	router.POST("/path", handlerFunc)

// 	// test not allowed
// 	// * (server)
// 	s := &fasthttp.Server{
// 		Handler: router.handler(),
// 	}

// 	rw := new(readWriter{})
// 	ch := make(chan error)

// 	rw.r.WriteString("OPTIONS * HTTP/1.1\r\nHost:\r\n\r\n")
// 	go func() {
// 		ch <- s.ServeConn(rw)
// 	}()
// 	select {
// 	case err := <-ch:
// 		if err != nil {
// 			t.Fatalf("return error %s", err)
// 		}
// 	case <-time.After(100 * time.Millisecond):
// 		t.Fatalf("timeout")
// 	}
// 	br := bufio.NewReader(&rw.w)
// 	var resp fasthttp.Response
// 	if err := resp.Read(br); err != nil {
// 		t.Fatalf("Unexpected error when reading response: %s", err)
// 	}
// 	if resp.Header.StatusCode() != fasthttp.StatusOK {
// 		t.Errorf("OPTIONS handling failed: Code=%d, Header=%v",
// 			resp.Header.StatusCode(), resp.Header.String())
// 	} else if allow := string(resp.Header.Peek("Allow")); allow != "POST, OPTIONS" {
// 		t.Error("unexpected Allow header value: " + allow)
// 	}

// 	// path
// 	rw.r.WriteString("OPTIONS /path HTTP/1.1\r\n\r\n")
// 	go func() {
// 		ch <- s.ServeConn(rw)
// 	}()
// 	select {
// 	case err := <-ch:
// 		if err != nil {
// 			t.Fatalf("return error %s", err)
// 		}
// 	case <-time.After(100 * time.Millisecond):
// 		t.Fatalf("timeout")
// 	}
// 	if err := resp.Read(br); err != nil {
// 		t.Fatalf("Unexpected error when reading response: %s", err)
// 	}
// 	if resp.Header.StatusCode() != fasthttp.StatusOK {
// 		t.Errorf("OPTIONS handling failed: Code=%d, Header=%v",
// 			resp.Header.StatusCode(), resp.Header.String())
// 	} else if allow := string(resp.Header.Peek("Allow")); allow != "POST, OPTIONS" {
// 		t.Error("unexpected Allow header value: " + allow)
// 	}

// 	rw.r.WriteString("OPTIONS /doesnotexist HTTP/1.1\r\n\r\n")
// 	go func() {
// 		ch <- s.ServeConn(rw)
// 	}()
// 	select {
// 	case err := <-ch:
// 		if err != nil {
// 			t.Fatalf("return error %s", err)
// 		}
// 	case <-time.After(100 * time.Millisecond):
// 		t.Fatalf("timeout")
// 	}
// 	if err := resp.Read(br); err != nil {
// 		t.Fatalf("Unexpected error when reading response: %s", err)
// 	}
// 	if !(resp.Header.StatusCode() == fasthttp.StatusNotFound) {
// 		t.Errorf("OPTIONS handling failed: Code=%d, Header=%v",
// 			resp.Header.StatusCode(), resp.Header.String())
// 	}

// 	// add another method
// 	router.GET("/path", handlerFunc)

// 	// test again
// 	// * (server)
// 	rw.r.WriteString("OPTIONS * HTTP/1.1\r\n\r\n")
// 	go func() {
// 		ch <- s.ServeConn(rw)
// 	}()
// 	select {
// 	case err := <-ch:
// 		if err != nil {
// 			t.Fatalf("return error %s", err)
// 		}
// 	case <-time.After(100 * time.Millisecond):
// 		t.Fatalf("timeout")
// 	}
// 	if err := resp.Read(br); err != nil {
// 		t.Fatalf("Unexpected error when reading response: %s", err)
// 	}
// 	if resp.Header.StatusCode() != fasthttp.StatusOK {
// 		t.Errorf("OPTIONS handling failed: Code=%d, Header=%v",
// 			resp.Header.StatusCode(), resp.Header.String())
// 	} else if allow := string(resp.Header.Peek("Allow")); allow != "POST, GET, OPTIONS" && allow != "GET, POST, OPTIONS" {
// 		t.Error("unexpected Allow header value: " + allow)
// 	}

// 	// path
// 	rw.r.WriteString("OPTIONS /path HTTP/1.1\r\n\r\n")
// 	go func() {
// 		ch <- s.ServeConn(rw)
// 	}()
// 	select {
// 	case err := <-ch:
// 		if err != nil {
// 			t.Fatalf("return error %s", err)
// 		}
// 	case <-time.After(100 * time.Millisecond):
// 		t.Fatalf("timeout")
// 	}
// 	if err := resp.Read(br); err != nil {
// 		t.Fatalf("Unexpected error when reading response: %s", err)
// 	}
// 	if resp.Header.StatusCode() != fasthttp.StatusOK {
// 		t.Errorf("OPTIONS handling failed: Code=%d, Header=%v",
// 			resp.Header.StatusCode(), resp.Header.String())
// 	} else if allow := string(resp.Header.Peek("Allow")); allow != "POST, GET, OPTIONS" && allow != "GET, POST, OPTIONS" {
// 		t.Error("unexpected Allow header value: " + allow)
// 	}

// 	// custom handler
// 	var custom bool
// 	router.OPTIONS("/path", func(_ *fasthttp.RequestCtx) {
// 		custom = true
// 	})

// 	// test again
// 	// * (server)
// 	rw.r.WriteString("OPTIONS * HTTP/1.1\r\n\r\n")
// 	go func() {
// 		ch <- s.ServeConn(rw)
// 	}()
// 	select {
// 	case err := <-ch:
// 		if err != nil {
// 			t.Fatalf("return error %s", err)
// 		}
// 	case <-time.After(100 * time.Millisecond):
// 		t.Fatalf("timeout")
// 	}
// 	if err := resp.Read(br); err != nil {
// 		t.Fatalf("Unexpected error when reading response: %s", err)
// 	}
// 	if resp.Header.StatusCode() != fasthttp.StatusOK {
// 		t.Errorf("OPTIONS handling failed: Code=%d, Header=%v",
// 			resp.Header.StatusCode(), resp.Header.String())
// 	} else if allow := string(resp.Header.Peek("Allow")); allow != "POST, GET, OPTIONS" && allow != "GET, POST, OPTIONS" {
// 		t.Error("unexpected Allow header value: " + allow)
// 	}
// 	if custom {
// 		t.Error("custom handler called on *")
// 	}

// 	// path
// 	rw.r.WriteString("OPTIONS /path HTTP/1.1\r\n\r\n")
// 	go func() {
// 		ch <- s.ServeConn(rw)
// 	}()
// 	select {
// 	case err := <-ch:
// 		if err != nil {
// 			t.Fatalf("return error %s", err)
// 		}
// 	case <-time.After(100 * time.Millisecond):
// 		t.Fatalf("timeout")
// 	}
// 	if err := resp.Read(br); err != nil {
// 		t.Fatalf("Unexpected error when reading response: %s", err)
// 	}
// 	if resp.Header.StatusCode() != fasthttp.StatusOK {
// 		t.Errorf("OPTIONS handling failed: Code=%d, Header=%v",
// 			resp.Header.StatusCode(), resp.Header.String())
// 	}
// 	if !custom {
// 		t.Error("custom handler not called")
// 	}
// }

func TestRouterNotAllowed(t *testing.T) {
	handlerFunc := func(_ *fasthttp.RequestCtx) {}

	router := New()
	router.POST("/path", handlerFunc)

	// Test not allowed
	s := &fasthttp.Server{
		Handler: router.handler(),
	}

	rw := new(readWriter)
	ch := make(chan error)

	rw.r.WriteString("GET /path HTTP/1.1\r\n\r\n")
	go func() {
		ch <- s.ServeConn(rw)
	}()
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("return error %s", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout")
	}
	br := bufio.NewReader(&rw.w)
	var resp fasthttp.Response
	if err := resp.Read(br); err != nil {
		t.Fatalf("Unexpected error when reading response: %s", err)
	}
	if !(resp.Header.StatusCode() == fasthttp.StatusMethodNotAllowed) {
		t.Errorf("NotAllowed handling failed: Code=%d. Actual=%d", resp.Header.StatusCode(), fasthttp.StatusMethodNotAllowed)
	} else if allow := string(resp.Header.Peek("Allow")); allow != "POST, OPTIONS" {
		t.Error("unexpected Allow header value: " + allow)
	}

	// add another method
	router.DELETE("/path", handlerFunc)
	router.OPTIONS("/path", handlerFunc) // must be ignored

	// test again
	rw.r.WriteString("GET /path HTTP/1.1\r\n\r\n")
	go func() {
		ch <- s.ServeConn(rw)
	}()
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("return error %s", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout")
	}
	if err := resp.Read(br); err != nil {
		t.Fatalf("Unexpected error when reading response: %s", err)
	}
	if !(resp.Header.StatusCode() == fasthttp.StatusMethodNotAllowed) {
		t.Errorf("NotAllowed handling failed: Code=%d. Actual=%d", resp.Header.StatusCode(), fasthttp.StatusMethodNotAllowed)
	} else if allow := string(resp.Header.Peek("Allow")); allow != "POST, DELETE, OPTIONS" && allow != "DELETE, POST, OPTIONS" {
		t.Error("unexpected Allow header value: " + allow)
	}

	responseText := "custom method"
	router.MethodNotAllowed(func(ctx *Context) {
		ctx.SetStatusCode(fasthttp.StatusTeapot)
		_, e := ctx.Write([]byte(responseText))
		_ = e
	})
	rw.r.WriteString("GET /path HTTP/1.1\r\n\r\n")
	go func() {
		ch <- s.ServeConn(rw)
	}()
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("return error %s", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout")
	}
	if err := resp.Read(br); err != nil {
		t.Fatalf("Unexpected error when reading response: %s", err)
	}
	if !bytes.Equal(resp.Body(), []byte(responseText)) {
		t.Errorf("unexpected response got %q want %q", string(resp.Body()), responseText)
	}
	if resp.Header.StatusCode() != fasthttp.StatusTeapot {
		t.Errorf("unexpected response code %d want %d", resp.Header.StatusCode(), fasthttp.StatusTeapot)
	}
	if allow := string(resp.Header.Peek("Allow")); allow != "POST, DELETE, OPTIONS" && allow != "DELETE, POST, OPTIONS" {
		t.Error("unexpected Allow header value: " + allow)
	}
}

func TestRouterNotFound(t *testing.T) {
	handlerFunc := func(_ *fasthttp.RequestCtx) {}

	router := New()
	router.GET("/path", handlerFunc)
	router.GET("/dir/", handlerFunc)
	router.GET("/", handlerFunc)
	router.GET("/path/:user", handlerFunc)
	router.Sub("/abc").GET("/:user", handlerFunc)

	testRoutes := []struct {
		route    string
		code     int
		location string
	}{
		{"/path/", 301, "/path"},                       // TSR -/
		{"/dir", 200, ""},                              // TSR -/
		{"/", 200, ""},                                 // TSR +/
		{"/PATH", 301, "/path"},                        // Fixed Case
		{"/DIR", 301, "/dir"},                          // Fixed Case
		{"/PATH/", 301, "/path"},                       // Fixed Case -/
		{"/DIR/", 301, "/dir"},                         // Fixed Case +/
		{"/paTh/?name=foo", 301, "/path?name=foo"},     // Fixed Case With Params +/
		{"/paTh?name=foo", 301, "/path?name=foo"},      // Fixed Case With Params +/
		{"/../path", 200, ""},                          // CleanPath
		{"/nope", 404, ""},                             // NotFound
		{"/path/?name=foo", 301, "/path?name=foo"},     // TSR Case With Params
		{"/path/u/?name=foo", 301, "/path/u?name=foo"}, // Dynamic TSR -/
		{"/abc/u/?name=foo", 301, "/abc/u?name=foo"},   // Sub Dynamic TSR -/
		{"/AbC/u/?name=foo", 301, "/abc/u?name=foo"},   // Sub Dynamic Fixed Case -/
	}

	s := &fasthttp.Server{
		Handler: router.handler(),
	}

	rw := new(readWriter)
	br := bufio.NewReader(&rw.w)
	var resp fasthttp.Response
	ch := make(chan error)
	for _, tr := range testRoutes {
		t.Logf("testing %v, want %v code", tr.route, tr.code)
		rw.r.WriteString(fmt.Sprintf("GET %s HTTP/1.1\r\n\r\n", tr.route))
		go func() {
			ch <- s.ServeConn(rw)
		}()
		select {
		case err := <-ch:
			if err != nil {
				t.Fatalf("return error %s", err)
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("timeout")
		}
		if err := resp.Read(br); err != nil {
			t.Fatalf("Unexpected error when reading response: %s", err)
		}
		if !(resp.Header.StatusCode() == tr.code) {
			t.Errorf("NotFound handling route %s failed: Code=%d want=%d",
				tr.route, resp.Header.StatusCode(), tr.code)
		}
		respLocation := string(resp.Header.Peek("Location"))
		if tr.code == 301 && respLocation != tr.location {
			t.Errorf("Wrong location header %s failed: Location=%s want=%s",
				tr.route, respLocation, tr.location)
		}
	}
	t.Log("not found test")
	// Test custom not found handler
	var notFound bool
	router.NotFound(func(ctx *Context) {
		ctx.SetStatusCode(404)
		notFound = true
	})
	rw.r.WriteString("GET /nope HTTP/1.1\r\n\r\n")
	go func() {
		ch <- s.ServeConn(rw)
	}()
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("return error %s", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout")
	}
	if err := resp.Read(br); err != nil {
		t.Fatalf("Unexpected error when reading response: %s", err)
	}
	if !(resp.Header.StatusCode() == http.StatusNotFound && notFound) {
		t.Errorf(
			"Custom NotFound handler failed: Code=%d, Header=%v, url=/nope",
			resp.Header.StatusCode(),
			string(resp.Header.Peek("Location")),
		)
	}

	// Test other method than GET (want 307 instead of 301)
	router.PATCH("/path", handlerFunc)
	rw.r.WriteString("PATCH /path/ HTTP/1.1\r\n\r\n")
	go func() {
		ch <- s.ServeConn(rw)
	}()
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("return error %s", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout")
	}
	if err := resp.Read(br); err != nil {
		t.Fatalf("Unexpected error when reading response: %s", err)
	}
	if resp.Header.StatusCode() != 307 {
		t.Errorf("Custom NotFound handler failed: Code=%d, Header=%v, url=/path",
			resp.Header.StatusCode(),
			string(resp.Header.Peek("Location")),
		)
	}

	// Test special case where no node for the prefix "/" exists
	router = New()
	router.GET("/a", handlerFunc)
	s.Handler = router.handler()
	rw.r.WriteString("GET / HTTP/1.1\r\n\r\n")
	go func() {
		ch <- s.ServeConn(rw)
	}()
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("return error %s", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout")
	}
	if err := resp.Read(br); err != nil {
		t.Fatalf("Unexpected error when reading response: %s", err)
	}
	if !(resp.Header.StatusCode() == 404) {
		t.Errorf("NotFound handling route / failed: Code=%d, Header=%v, url=/",
			resp.Header.StatusCode(),
			string(resp.Header.Peek("Location")),
		)
	}
}

func TestRouterPanicHandler(t *testing.T) {
	router := New()
	panicHandled := false

	router.PanicHandler(func(ctx *Context, p interface{}) {
		panicHandled = true
	})

	router.Handle("PUT", "/user/:name", func(_ *fasthttp.RequestCtx) {
		panic("oops!")
	})

	defer func() {
		if rcv := recover(); rcv != nil {
			t.Fatal("handling panic failed")
		}
	}()

	s := &fasthttp.Server{
		Handler: router.handler(),
	}

	rw := new(readWriter)
	ch := make(chan error)

	rw.r.WriteString(string("PUT /user/gopher HTTP/1.1\r\n\r\n"))
	go func() {
		ch <- s.ServeConn(rw)
	}()
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("return error %s", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout")
	}

	if !panicHandled {
		t.Fatal("simulating failed")
	}
}

func TestRouterLookup(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			Logger.Errorf("panic handled: %v", r)
			debug.PrintStack()
		}
	}()
	routed := false
	wantHandle := func(_ *fasthttp.RequestCtx) {
		routed = true
	}

	router := New()
	ctx := &Context{
		RequestCtx: &fasthttp.RequestCtx{},
		Logger:     Logger,
		App:        router,
	}

	// try empty router first
	handle, tsr := router.defaultRouter.Lookup("GET", "/nope", ctx)
	if handle != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", handle)
	}
	if tsr {
		t.Error("Got wrong TSR recommendation!")
	}

	// insert route and try again
	router.GET("/user/:name", wantHandle)

	handle, _ = router.defaultRouter.Lookup("GET", "/user/gopher", ctx)
	if handle == nil {
		t.Fatal("Got no handle!")
	} else {
		handle(nil)
		if !routed {
			t.Fatal("Routing failed!")
		}
	}
	if ctx.UserValue("name") != "gopher" {
		t.Error("Param not set!")
	}

	handle, tsr = router.defaultRouter.Lookup("GET", "/user/gopher/", ctx)
	if handle != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", handle)
	}
	if !tsr {
		t.Error("Got no TSR recommendation!")
	}

	handle, tsr = router.defaultRouter.Lookup("GET", "/nope", ctx)
	if handle != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", handle)
	}
	if tsr {
		t.Error("Got wrong TSR recommendation!")
	}
}

// func (mfs *mockFileSystem) Open(name string) (http.File, error) {
// 	mfs.opened = true
// 	return nil, errors.New("this is just a mock")
// }

func TestRouterServeFiles(t *testing.T) {
	router := New()

	recv := catchPanic(func() {
		router.defaultRouter.ServeFiles("/noFilepath", os.TempDir())
	})
	if recv == nil {
		t.Fatal("registering path not ending with '*filepath' did not panic")
	}
	body := []byte("fake ico")
	err := os.WriteFile(os.TempDir()+"/favicon.ico", body, 0644)
	if err != nil {
		t.Fatal(err)
	}

	router.defaultRouter.ServeFiles("/*filepath", os.TempDir())

	s := &fasthttp.Server{
		Handler: router.handler(),
	}

	rw := new(readWriter)
	ch := make(chan error)

	rw.r.WriteString(string("GET /favicon.ico HTTP/1.1\r\n\r\n"))
	go func() {
		ch <- s.ServeConn(rw)
	}()
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("return error %s", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("timeout")
	}

	br := bufio.NewReader(&rw.w)
	var resp fasthttp.Response
	if err := resp.Read(br); err != nil {
		t.Fatalf("Unexpected error when reading response: %s", err)
	}
	if resp.Header.StatusCode() != 200 {
		t.Fatalf("Unexpected status code %d. Expected %d", resp.Header.StatusCode(), 423)
	}
	if !bytes.Equal(resp.Body(), body) {
		t.Fatalf("Unexpected body %q. Expected %q", resp.Body(), string(body))
	}
}

func (rw *readWriter) Close() error {
	return nil
}

func (rw *readWriter) Read(b []byte) (int, error) {
	return rw.r.Read(b)
}

func (rw *readWriter) Write(b []byte) (int, error) {
	return rw.w.Write(b)
}

func (rw *readWriter) RemoteAddr() net.Addr {
	return zeroTCPAddr
}

func (rw *readWriter) LocalAddr() net.Addr {
	return zeroTCPAddr
}

func (rw *readWriter) SetReadDeadline(t time.Time) error {
	return nil
}

func (rw *readWriter) SetWriteDeadline(t time.Time) error {
	return nil
}
