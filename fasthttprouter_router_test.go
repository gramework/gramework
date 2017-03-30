// Copyright 2013 Julien Schmidt. All rights reserved.
// Copyright (c) 2015-2016, 招牌疯子
// Copyright (c) 2017, Kirill Danshin
// Use of this source code is governed by a BSD-style license that can be found
// in the 3rd-Party License/fasthttprouter file.

package gramework

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
)

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

	rw := &readWriter{}
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

type handlerStruct struct {
	handeled *bool
}

func (h handlerStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	*h.handeled = true
}

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

	rw := &readWriter{}
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

// 	rw := &readWriter{}
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

	rw := &readWriter{}
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
		t.Errorf("NotAllowed handling failed: Code=%d", resp.Header.StatusCode())
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
		t.Errorf("NotAllowed handling failed: Code=%d", resp.Header.StatusCode())
	} else if allow := string(resp.Header.Peek("Allow")); allow != "POST, DELETE, OPTIONS" && allow != "DELETE, POST, OPTIONS" {
		t.Error("unexpected Allow header value: " + allow)
	}

	responseText := "custom method"
	router.MethodNotAllowed(func(ctx *Context) {
		ctx.SetStatusCode(fasthttp.StatusTeapot)
		ctx.Write([]byte(responseText))
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

	testRoutes := []struct {
		route string
		code  int
	}{
		{"/path/", 301},          // TSR -/
		{"/dir", 301},            // TSR +/
		{"/", 200},               // TSR +/
		{"/PATH", 301},           // Fixed Case
		{"/DIR", 301},            // Fixed Case
		{"/PATH/", 301},          // Fixed Case -/
		{"/DIR/", 301},           // Fixed Case +/
		{"/paTh/?name=foo", 301}, // Fixed Case With Params +/
		{"/paTh?name=foo", 301},  // Fixed Case With Params +/
		{"/../path", 200},        // CleanPath
		{"/nope", 404},           // NotFound
	}

	s := &fasthttp.Server{
		Handler: router.handler(),
	}

	rw := &readWriter{}
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
	if !(resp.Header.StatusCode() == 404 && notFound == true) {
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

	rw := &readWriter{}
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

	handle, tsr = router.defaultRouter.Lookup("GET", "/user/gopher", ctx)
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

type mockFileSystem struct {
	opened bool
}

func (mfs *mockFileSystem) Open(name string) (http.File, error) {
	mfs.opened = true
	return nil, errors.New("this is just a mock")
}

func TestRouterServeFiles(t *testing.T) {
	router := New()

	recv := catchPanic(func() {
		router.defaultRouter.ServeFiles("/noFilepath", os.TempDir())
	})
	if recv == nil {
		t.Fatal("registering path not ending with '*filepath' did not panic")
	}
	body := []byte("fake ico")
	ioutil.WriteFile(os.TempDir()+"/favicon.ico", body, 0644)

	router.defaultRouter.ServeFiles("/*filepath", os.TempDir())

	s := &fasthttp.Server{
		Handler: router.handler(),
	}

	rw := &readWriter{}
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

type readWriter struct {
	net.Conn
	r bytes.Buffer
	w bytes.Buffer
}

var zeroTCPAddr = &net.TCPAddr{
	IP: net.IPv4zero,
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
