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
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/gramework/gramework/x/testutils"

	"github.com/valyala/fasthttp"
)

func TestGrameRouter(t *testing.T) {
	app := New()
	if h, _ := app.defaultRouter.Lookup("GET", "/", nil); h != nil {
		t.Log("GET / should not return handler before registration")
		t.FailNow()
	}
	app.GET("/", 1)
	if h, _ := app.defaultRouter.Lookup("GET", "/", nil); h == nil {
		t.Log("GET / should return handler after registration")
		t.FailNow()
	}
	if h, _ := app.defaultRouter.Lookup("GET", "/abc", nil); h != nil {
		t.Log("GET /abc should not return handler before registration")
		t.FailNow()
	}

	app.GET("/abc", "abc")

	if h, _ := app.defaultRouter.Lookup("GET", "/abc", nil); h == nil {
		t.Log("GET /abc should return handler after registration")
		t.FailNow()
	}

	app.GET("/redir", app.ToTLSHandler())

	// POST

	if h, _ := app.defaultRouter.Lookup("POST", "/", nil); h != nil {
		fpc := runtime.FuncForPC(reflect.ValueOf(h).Pointer())
		file, line := fpc.FileLine(fpc.Entry())
		t.Logf("POST / should not return handler before registration, got %q (%v:%v)", fpc.Name(), file, line)
		t.FailNow()
	}
	app.POST("/", 1)
	if h, _ := app.defaultRouter.Lookup("POST", "/", nil); h == nil {
		fpc := runtime.FuncForPC(reflect.ValueOf(h).Pointer())
		file, line := fpc.FileLine(fpc.Entry())
		t.Logf("POST / should return handler after registration, got %q (%v:%v)", fpc.Name(), file, line)
		t.FailNow()
	}
	if h, _ := app.defaultRouter.Lookup("POST", "/abc", nil); h != nil {
		fpc := runtime.FuncForPC(reflect.ValueOf(h).Pointer())
		file, line := fpc.FileLine(fpc.Entry())
		t.Logf("POST /abc should not return handler before registration, got %q (%v:%v)", fpc.Name(), file, line)
		t.FailNow()
	}

	app.POST("/abc", "abc")

	if h, _ := app.defaultRouter.Lookup("POST", "/abc", nil); h == nil {
		t.Log("POST /abc should return handler after registration")
		t.FailNow()
	}

	// PUT
	if h, _ := app.defaultRouter.Lookup("PUT", "/", nil); h != nil {
		t.Log("PUT / should not return handler before registration")
		t.FailNow()
	}
	app.PUT("/", 1)
	if h, _ := app.defaultRouter.Lookup("PUT", "/", nil); h == nil {
		t.Log("PUT / should return handler after registration")
		t.FailNow()
	}
	if h, _ := app.defaultRouter.Lookup("PUT", "/abc", nil); h != nil {
		t.Log("PUT /abc should not return handler before registration")
		t.FailNow()
	}

	app.PUT("/abc", "abc")

	if h, _ := app.defaultRouter.Lookup("PUT", "/abc", nil); h == nil {
		t.Log("PUT /abc should return handler after registration")
		t.FailNow()
	}

	// DELETE
	if h, _ := app.defaultRouter.Lookup("DELETE", "/", nil); h != nil {
		t.Log("DELETE / should not return handler before registration")
		t.FailNow()
	}
	app.DELETE("/", 1)
	if h, _ := app.defaultRouter.Lookup("DELETE", "/", nil); h == nil {
		t.Log("DELETE / should return handler after registration")
		t.FailNow()
	}
	if h, _ := app.defaultRouter.Lookup("DELETE", "/abc", nil); h != nil {
		t.Log("DELETE /abc should not return handler before registration")
		t.FailNow()
	}

	app.DELETE("/abc", "abc")

	if h, _ := app.defaultRouter.Lookup("DELETE", "/abc", nil); h == nil {
		t.Log("DELETE /abc should return handler after registration")
		t.FailNow()
	}

	// HEAD
	if h, _ := app.defaultRouter.Lookup("HEAD", "/", nil); h != nil {
		t.Log("HEAD / should not return handler before registration")
		t.FailNow()
	}
	app.HEAD("/", 1)
	if h, _ := app.defaultRouter.Lookup("HEAD", "/", nil); h == nil {
		t.Log("HEAD / should return handler after registration")
		t.FailNow()
	}
	if h, _ := app.defaultRouter.Lookup("HEAD", "/abc", nil); h != nil {
		t.Log("HEAD /abc should not return handler before registration")
		t.FailNow()
	}

	app.HEAD("/abc", "abc")

	if h, _ := app.defaultRouter.Lookup("HEAD", "/abc", nil); h == nil {
		t.Log("HEAD /abc should return handler after registration")
		t.FailNow()
	}

	// OPTIONS
	if h, _ := app.defaultRouter.Lookup("OPTIONS", "/", nil); h != nil {
		t.Log("OPTIONS / should not return handler before registration")
		t.FailNow()
	}
	app.OPTIONS("/", 1)
	if h, _ := app.defaultRouter.Lookup("OPTIONS", "/", nil); h == nil {
		t.Log("OPTIONS / should return handler after registration")
		t.FailNow()
	}
	if h, _ := app.defaultRouter.Lookup("OPTIONS", "/abc", nil); h != nil {
		t.Log("OPTIONS /abc should not return handler before registration")
		t.FailNow()
	}

	app.OPTIONS("/abc", "abc")

	if h, _ := app.defaultRouter.Lookup("OPTIONS", "/abc", nil); h == nil {
		t.Log("OPTIONS /abc should return handler after registration")
		t.FailNow()
	}

	// PATCH
	if h, _ := app.defaultRouter.Lookup("PATCH", "/", nil); h != nil {
		t.Log("PATCH / should not return handler before registration")
		t.FailNow()
	}
	app.PATCH("/", 1)
	if h, _ := app.defaultRouter.Lookup("PATCH", "/", nil); h == nil {
		t.Log("PATCH / should return handler after registration")
		t.FailNow()
	}
	if h, _ := app.defaultRouter.Lookup("PATCH", "/abc", nil); h != nil {
		t.Log("PATCH /abc should not return handler before registration")
		t.FailNow()
	}

	app.PATCH("/abc", "abc")

	if h, _ := app.defaultRouter.Lookup("PATCH", "/abc", nil); h == nil {
		t.Log("PATCH /abc should return handler after registration")
		t.FailNow()
	}

	app.Sub("/abc").Handle("GET", "/def", "abcdef")

	if h, _ := app.defaultRouter.Lookup("GET", "/abc/def", nil); h == nil {
		t.Log("GET /abc/def should return handler after registration")
		t.FailNow()
	}

	app.Handle("CONNECT", "/ws", "")
	app.PanicHandler(nil)
	app.NotFound(nil)
	app.HandleMethodNotAllowed(true)
	app.HandleOPTIONS(true)

	port := testutils.Port().NonRoot().Unused().Acquire()
	bindAddr := fmt.Sprintf(":%d", port)
	go func() {
		err := app.ListenAndServe(bindAddr)
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(250 * time.Millisecond)
	_, err := http.Get("http://127.0.0.1" + bindAddr) // just should not panic
	if err != nil {
		t.Error(err)
	}
}

func TestSubRouter(t *testing.T) {
	app := New()

	app.Sub("/abc").Handle("GET", "/def", "abcdef")
	if h, _ := app.defaultRouter.Lookup("GET", "/abc/def", nil); h == nil {
		t.Log("GET /abc/def should return handler after registration")
		t.FailNow()
	}

	app.Sub("/abc").Handle("GET", "/def2/", "abcdef")
	if h, _ := app.defaultRouter.Lookup("GET", "/abc/def2", nil); h == nil {
		t.Log("GET /abc/def should return handler after registration")
		t.FailNow()
	}

	app.Sub("/cba").Handle("GET", "/def/:user", "abcdef")
	if h, _ := app.defaultRouter.Lookup("GET", "/cba/def/usr", nil); h == nil {
		t.Log("GET /cba/def/usr should return handler after registration")
		t.FailNow()
	}

	app.Sub("/hello").Sub("/world").Handle("POST", "/def", "defdef")
	if h, _ := app.defaultRouter.Lookup("POST", "/hello/world/def", nil); h == nil {
		t.Log("GET /hello/world/def/usr should return handler after registration")
		t.FailNow()
	}

	app.Sub("/hello").Sub("/world").Handle("POST", "/def/:user", "user")
	if h, _ := app.defaultRouter.Lookup("POST", "/hello/world/def/usr", nil); h == nil {
		t.Log("GET /hello/world/def/usr should return handler after registration")
		t.FailNow()
	}
}

func TestDomainRouter(t *testing.T) {
	app := New()
	r := app.Domain("example.com")
	if h, _ := r.Lookup("GET", "/", nil); h != nil {
		t.Log("GET / should not return handler before registration")
		t.FailNow()
	}
	r.GET("/", 1)
	if h, _ := r.Lookup("GET", "/", nil); h == nil {
		t.Log("GET / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("GET", "/abc", nil); h != nil {
		t.Log("GET /abc should not return handler before registration")
		t.FailNow()
	}

	r.GET("/abc", "abc")

	if h, _ := r.Lookup("GET", "/abc", nil); h == nil {
		t.Log("GET /abc should return handler after registration")
		t.FailNow()
	}

	r.GET("/redir", app.ToTLSHandler())

	// POST

	if h, _ := r.Lookup("POST", "/", nil); h != nil {
		t.Log("POST / should not return handler before registration")
		t.FailNow()
	}
	r.POST("/", 1)
	if h, _ := r.Lookup("POST", "/", nil); h == nil {
		t.Log("POST / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("POST", "/abc", nil); h != nil {
		t.Log("POST /abc should not return handler before registration")
		t.FailNow()
	}

	r.POST("/abc", "abc")

	if h, _ := r.Lookup("POST", "/abc", nil); h == nil {
		t.Log("POST /abc should return handler after registration")
		t.FailNow()
	}

	// PUT
	if h, _ := r.Lookup("PUT", "/", nil); h != nil {
		t.Log("PUT / should not return handler before registration")
		t.FailNow()
	}
	r.PUT("/", 1)
	if h, _ := r.Lookup("PUT", "/", nil); h == nil {
		t.Log("PUT / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("PUT", "/abc", nil); h != nil {
		t.Log("PUT /abc should not return handler before registration")
		t.FailNow()
	}

	r.PUT("/abc", "abc")

	if h, _ := r.Lookup("PUT", "/abc", nil); h == nil {
		t.Log("PUT /abc should return handler after registration")
		t.FailNow()
	}

	// DELETE
	if h, _ := r.Lookup("DELETE", "/", nil); h != nil {
		t.Log("DELETE / should not return handler before registration")
		t.FailNow()
	}
	r.DELETE("/", 1)
	if h, _ := r.Lookup("DELETE", "/", nil); h == nil {
		t.Log("DELETE / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("DELETE", "/abc", nil); h != nil {
		t.Log("DELETE /abc should not return handler before registration")
		t.FailNow()
	}

	r.DELETE("/abc", "abc")

	if h, _ := r.Lookup("DELETE", "/abc", nil); h == nil {
		t.Log("DELETE /abc should return handler after registration")
		t.FailNow()
	}

	// HEAD
	if h, _ := r.Lookup("HEAD", "/", nil); h != nil {
		t.Log("HEAD / should not return handler before registration")
		t.FailNow()
	}
	r.HEAD("/", 1)
	if h, _ := r.Lookup("HEAD", "/", nil); h == nil {
		t.Log("HEAD / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("HEAD", "/abc", nil); h != nil {
		t.Log("HEAD /abc should not return handler before registration")
		t.FailNow()
	}

	r.HEAD("/abc", "abc")

	if h, _ := r.Lookup("HEAD", "/abc", nil); h == nil {
		t.Log("HEAD /abc should return handler after registration")
		t.FailNow()
	}

	// OPTIONS
	if h, _ := r.Lookup("OPTIONS", "/", nil); h != nil {
		t.Log("OPTIONS / should not return handler before registration")
		t.FailNow()
	}
	r.OPTIONS("/", 1)
	if h, _ := r.Lookup("OPTIONS", "/", nil); h == nil {
		t.Log("OPTIONS / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("OPTIONS", "/abc", nil); h != nil {
		t.Log("OPTIONS /abc should not return handler before registration")
		t.FailNow()
	}

	r.OPTIONS("/abc", "abc")

	if h, _ := r.Lookup("OPTIONS", "/abc", nil); h == nil {
		t.Log("OPTIONS /abc should return handler after registration")
		t.FailNow()
	}

	// PATCH
	if h, _ := r.Lookup("PATCH", "/", nil); h != nil {
		t.Log("PATCH / should not return handler before registration")
		t.FailNow()
	}
	r.PATCH("/", 1)
	if h, _ := r.Lookup("PATCH", "/", nil); h == nil {
		t.Log("PATCH / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("PATCH", "/abc", nil); h != nil {
		t.Log("PATCH /abc should not return handler before registration")
		t.FailNow()
	}

	r.PATCH("/abc", "abc")

	if h, _ := r.Lookup("PATCH", "/abc", nil); h == nil {
		t.Log("PATCH /abc should return handler after registration")
		t.FailNow()
	}

	port := testutils.Port().NonRoot().Unused().Acquire()
	bindAddr := fmt.Sprintf(":%d", port)
	go func() {
		err := app.ListenAndServe(bindAddr)
		if err != nil {
			panic(err)
		}
	}()
	_, err := http.Get("http://127.0.0.1" + bindAddr) // just should not panic
	_ = err
}

func TestDomainHTTPRouter(t *testing.T) {
	app := New()
	r := app.Domain("example.com").HTTP()
	if h, _ := r.Lookup("GET", "/", nil); h != nil {
		t.Log("GET / should not return handler before registration")
		t.FailNow()
	}
	r.GET("/", 1)
	if h, _ := r.Lookup("GET", "/", nil); h == nil {
		t.Log("GET / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("GET", "/abc", nil); h != nil {
		t.Log("GET /abc should not return handler before registration")
		t.FailNow()
	}

	r.GET("/abc", "abc")

	if h, _ := r.Lookup("GET", "/abc", nil); h == nil {
		t.Log("GET /abc should return handler after registration")
		t.FailNow()
	}

	r.GET("/redir", app.ToTLSHandler())

	// POST

	if h, _ := r.Lookup("POST", "/", nil); h != nil {
		t.Log("POST / should not return handler before registration")
		t.FailNow()
	}
	r.POST("/", 1)
	if h, _ := r.Lookup("POST", "/", nil); h == nil {
		t.Log("POST / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("POST", "/abc", nil); h != nil {
		t.Log("POST /abc should not return handler before registration")
		t.FailNow()
	}

	r.POST("/abc", "abc")

	if h, _ := r.Lookup("POST", "/abc", nil); h == nil {
		t.Log("POST /abc should return handler after registration")
		t.FailNow()
	}

	// PUT
	if h, _ := r.Lookup("PUT", "/", nil); h != nil {
		t.Log("PUT / should not return handler before registration")
		t.FailNow()
	}
	r.PUT("/", 1)
	if h, _ := r.Lookup("PUT", "/", nil); h == nil {
		t.Log("PUT / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("PUT", "/abc", nil); h != nil {
		t.Log("PUT /abc should not return handler before registration")
		t.FailNow()
	}

	r.PUT("/abc", "abc")

	if h, _ := r.Lookup("PUT", "/abc", nil); h == nil {
		t.Log("PUT /abc should return handler after registration")
		t.FailNow()
	}

	// DELETE
	if h, _ := r.Lookup("DELETE", "/", nil); h != nil {
		t.Log("DELETE / should not return handler before registration")
		t.FailNow()
	}
	r.DELETE("/", 1)
	if h, _ := r.Lookup("DELETE", "/", nil); h == nil {
		t.Log("DELETE / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("DELETE", "/abc", nil); h != nil {
		t.Log("DELETE /abc should not return handler before registration")
		t.FailNow()
	}

	r.DELETE("/abc", "abc")

	if h, _ := r.Lookup("DELETE", "/abc", nil); h == nil {
		t.Log("DELETE /abc should return handler after registration")
		t.FailNow()
	}

	// HEAD
	if h, _ := r.Lookup("HEAD", "/", nil); h != nil {
		t.Log("HEAD / should not return handler before registration")
		t.FailNow()
	}
	r.HEAD("/", 1)
	if h, _ := r.Lookup("HEAD", "/", nil); h == nil {
		t.Log("HEAD / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("HEAD", "/abc", nil); h != nil {
		t.Log("HEAD /abc should not return handler before registration")
		t.FailNow()
	}

	r.HEAD("/abc", "abc")

	if h, _ := r.Lookup("HEAD", "/abc", nil); h == nil {
		t.Log("HEAD /abc should return handler after registration")
		t.FailNow()
	}

	// OPTIONS
	if h, _ := r.Lookup("OPTIONS", "/", nil); h != nil {
		t.Log("OPTIONS / should not return handler before registration")
		t.FailNow()
	}
	r.OPTIONS("/", 1)
	if h, _ := r.Lookup("OPTIONS", "/", nil); h == nil {
		t.Log("OPTIONS / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("OPTIONS", "/abc", nil); h != nil {
		t.Log("OPTIONS /abc should not return handler before registration")
		t.FailNow()
	}

	r.OPTIONS("/abc", "abc")

	if h, _ := r.Lookup("OPTIONS", "/abc", nil); h == nil {
		t.Log("OPTIONS /abc should return handler after registration")
		t.FailNow()
	}

	// PATCH
	if h, _ := r.Lookup("PATCH", "/", nil); h != nil {
		t.Log("PATCH / should not return handler before registration")
		t.FailNow()
	}
	r.PATCH("/", 1)
	if h, _ := r.Lookup("PATCH", "/", nil); h == nil {
		t.Log("PATCH / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("PATCH", "/abc", nil); h != nil {
		t.Log("PATCH /abc should not return handler before registration")
		t.FailNow()
	}

	r.PATCH("/abc", "abc")

	if h, _ := r.Lookup("PATCH", "/abc", nil); h == nil {
		t.Log("PATCH /abc should return handler after registration")
		t.FailNow()
	}
}

func TestDomainHTTPSRouter(t *testing.T) {
	app := New()
	r := app.Domain("example.com").HTTPS()
	if h, _ := r.Lookup("GET", "/", nil); h != nil {
		t.Log("GET / should not return handler before registration")
		t.FailNow()
	}
	r.GET("/", 1)
	if h, _ := r.Lookup("GET", "/", nil); h == nil {
		t.Log("GET / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("GET", "/abc", nil); h != nil {
		t.Log("GET /abc should not return handler before registration")
		t.FailNow()
	}

	r.GET("/abc", "abc")

	if h, _ := r.Lookup("GET", "/abc", nil); h == nil {
		t.Log("GET /abc should return handler after registration")
		t.FailNow()
	}

	r.GET("/redir", app.ToTLSHandler())

	if h, _ := r.Lookup("GET", "/redir", nil); h == nil {
		t.Log("GET /abc should return handler after registration")
		t.FailNow()
	} else {
		defer func() {
			e := recover()
			if e != nil {
				t.Log("panic handled when testing /redir")
			}
		}()
		h(&Context{
			RequestCtx: &fasthttp.RequestCtx{},
		})
	}

	// POST

	if h, _ := r.Lookup("POST", "/", nil); h != nil {
		t.Log("POST / should not return handler before registration")
		t.FailNow()
	}
	r.POST("/", 1)
	if h, _ := r.Lookup("POST", "/", nil); h == nil {
		t.Log("POST / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("POST", "/abc", nil); h != nil {
		t.Log("POST /abc should not return handler before registration")
		t.FailNow()
	}

	r.POST("/abc", "abc")

	if h, _ := r.Lookup("POST", "/abc", nil); h == nil {
		t.Log("POST /abc should return handler after registration")
		t.FailNow()
	}

	// PUT
	if h, _ := r.Lookup("PUT", "/", nil); h != nil {
		t.Log("PUT / should not return handler before registration")
		t.FailNow()
	}
	r.PUT("/", 1)
	if h, _ := r.Lookup("PUT", "/", nil); h == nil {
		t.Log("PUT / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("PUT", "/abc", nil); h != nil {
		t.Log("PUT /abc should not return handler before registration")
		t.FailNow()
	}

	r.PUT("/abc", "abc")

	if h, _ := r.Lookup("PUT", "/abc", nil); h == nil {
		t.Log("PUT /abc should return handler after registration")
		t.FailNow()
	}

	// DELETE
	if h, _ := r.Lookup("DELETE", "/", nil); h != nil {
		t.Log("DELETE / should not return handler before registration")
		t.FailNow()
	}
	r.DELETE("/", 1)
	if h, _ := r.Lookup("DELETE", "/", nil); h == nil {
		t.Log("DELETE / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("DELETE", "/abc", nil); h != nil {
		t.Log("DELETE /abc should not return handler before registration")
		t.FailNow()
	}

	r.DELETE("/abc", "abc")

	if h, _ := r.Lookup("DELETE", "/abc", nil); h == nil {
		t.Log("DELETE /abc should return handler after registration")
		t.FailNow()
	}

	// HEAD
	if h, _ := r.Lookup("HEAD", "/", nil); h != nil {
		t.Log("HEAD / should not return handler before registration")
		t.FailNow()
	}
	r.HEAD("/", 1)
	if h, _ := r.Lookup("HEAD", "/", nil); h == nil {
		t.Log("HEAD / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("HEAD", "/abc", nil); h != nil {
		t.Log("HEAD /abc should not return handler before registration")
		t.FailNow()
	}

	r.HEAD("/abc", "abc")

	if h, _ := r.Lookup("HEAD", "/abc", nil); h == nil {
		t.Log("HEAD /abc should return handler after registration")
		t.FailNow()
	}

	// OPTIONS
	if h, _ := r.Lookup("OPTIONS", "/", nil); h != nil {
		t.Log("OPTIONS / should not return handler before registration")
		t.FailNow()
	}
	r.OPTIONS("/", 1)
	if h, _ := r.Lookup("OPTIONS", "/", nil); h == nil {
		t.Log("OPTIONS / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("OPTIONS", "/abc", nil); h != nil {
		t.Log("OPTIONS /abc should not return handler before registration")
		t.FailNow()
	}

	r.OPTIONS("/abc", "abc")

	if h, _ := r.Lookup("OPTIONS", "/abc", nil); h == nil {
		t.Log("OPTIONS /abc should return handler after registration")
		t.FailNow()
	}

	// PATCH
	if h, _ := r.Lookup("PATCH", "/", nil); h != nil {
		t.Log("PATCH / should not return handler before registration")
		t.FailNow()
	}
	r.PATCH("/", 1)
	if h, _ := r.Lookup("PATCH", "/", nil); h == nil {
		t.Log("PATCH / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("PATCH", "/abc", nil); h != nil {
		t.Log("PATCH /abc should not return handler before registration")
		t.FailNow()
	}

	r.PATCH("/abc", "abc")

	if h, _ := r.Lookup("PATCH", "/abc", nil); h == nil {
		t.Log("PATCH /abc should return handler after registration")
		t.FailNow()
	}

	r.GET("/err", func(ctx *Context) error {
		return errors.New("test")
	})

	r.GET("/fasterr", func(ctx *fasthttp.RequestCtx) error {
		return errors.New("test")
	})

	if h, _ := r.Lookup("PATCH", "/abc", nil); h == nil {
		t.Log("PATCH /abc should return handler after registration")
		t.FailNow()
	}

	port := testutils.Port().NonRoot().Unused().Acquire()
	bindAddr := fmt.Sprintf(":%d", port)
	go func() {
		err := app.ListenAndServe(bindAddr)
		if err != nil {
			panic(err)
		}
	}()
	_, e := http.Get("http://127.0.0.1" + bindAddr) // just should not panic
	_ = e
	_, e = http.Get("http://127.0.0.1" + bindAddr) // just should not panic, twice
	_ = e
}

func TestHTTPRouter(t *testing.T) {
	app := New()
	r := app.HTTP()
	if h, _ := r.Lookup("GET", "/", nil); h != nil {
		t.Log("GET / should not return handler before registration")
		t.FailNow()
	}
	r.GET("/", 1)
	if h, _ := r.Lookup("GET", "/", nil); h == nil {
		t.Log("GET / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("GET", "/abc", nil); h != nil {
		t.Log("GET /abc should not return handler before registration")
		t.FailNow()
	}

	r.GET("/abc", "abc")

	if h, _ := r.Lookup("GET", "/abc", nil); h == nil {
		t.Log("GET /abc should return handler after registration")
		t.FailNow()
	}

	r.GET("/redir", app.ToTLSHandler())

	// POST

	if h, _ := r.Lookup("POST", "/", nil); h != nil {
		t.Log("POST / should not return handler before registration")
		t.FailNow()
	}
	r.POST("/", 1)
	if h, _ := r.Lookup("POST", "/", nil); h == nil {
		t.Log("POST / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("POST", "/abc", nil); h != nil {
		t.Log("POST /abc should not return handler before registration")
		t.FailNow()
	}

	r.POST("/abc", "abc")

	if h, _ := r.Lookup("POST", "/abc", nil); h == nil {
		t.Log("POST /abc should return handler after registration")
		t.FailNow()
	}

	// PUT
	if h, _ := r.Lookup("PUT", "/", nil); h != nil {
		t.Log("PUT / should not return handler before registration")
		t.FailNow()
	}
	r.PUT("/", 1)
	if h, _ := r.Lookup("PUT", "/", nil); h == nil {
		t.Log("PUT / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("PUT", "/abc", nil); h != nil {
		t.Log("PUT /abc should not return handler before registration")
		t.FailNow()
	}

	r.PUT("/abc", "abc")

	if h, _ := r.Lookup("PUT", "/abc", nil); h == nil {
		t.Log("PUT /abc should return handler after registration")
		t.FailNow()
	}

	// DELETE
	if h, _ := r.Lookup("DELETE", "/", nil); h != nil {
		t.Log("DELETE / should not return handler before registration")
		t.FailNow()
	}
	r.DELETE("/", 1)
	if h, _ := r.Lookup("DELETE", "/", nil); h == nil {
		t.Log("DELETE / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("DELETE", "/abc", nil); h != nil {
		t.Log("DELETE /abc should not return handler before registration")
		t.FailNow()
	}

	r.DELETE("/abc", "abc")

	if h, _ := r.Lookup("DELETE", "/abc", nil); h == nil {
		t.Log("DELETE /abc should return handler after registration")
		t.FailNow()
	}

	// HEAD
	if h, _ := r.Lookup("HEAD", "/", nil); h != nil {
		t.Log("HEAD / should not return handler before registration")
		t.FailNow()
	}
	r.HEAD("/", 1)
	if h, _ := r.Lookup("HEAD", "/", nil); h == nil {
		t.Log("HEAD / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("HEAD", "/abc", nil); h != nil {
		t.Log("HEAD /abc should not return handler before registration")
		t.FailNow()
	}

	r.HEAD("/abc", "abc")

	if h, _ := r.Lookup("HEAD", "/abc", nil); h == nil {
		t.Log("HEAD /abc should return handler after registration")
		t.FailNow()
	}

	// OPTIONS
	if h, _ := r.Lookup("OPTIONS", "/", nil); h != nil {
		t.Log("OPTIONS / should not return handler before registration")
		t.FailNow()
	}
	r.OPTIONS("/", 1)
	if h, _ := r.Lookup("OPTIONS", "/", nil); h == nil {
		t.Log("OPTIONS / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("OPTIONS", "/abc", nil); h != nil {
		t.Log("OPTIONS /abc should not return handler before registration")
		t.FailNow()
	}

	r.OPTIONS("/abc", "abc")

	if h, _ := r.Lookup("OPTIONS", "/abc", nil); h == nil {
		t.Log("OPTIONS /abc should return handler after registration")
		t.FailNow()
	}

	// PATCH
	if h, _ := r.Lookup("PATCH", "/", nil); h != nil {
		t.Log("PATCH / should not return handler before registration")
		t.FailNow()
	}
	r.PATCH("/", 1)
	if h, _ := r.Lookup("PATCH", "/", nil); h == nil {
		t.Log("PATCH / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("PATCH", "/abc", nil); h != nil {
		t.Log("PATCH /abc should not return handler before registration")
		t.FailNow()
	}

	r.PATCH("/abc", "abc")

	if h, _ := r.Lookup("PATCH", "/abc", nil); h == nil {
		t.Log("PATCH /abc should return handler after registration")
		t.FailNow()
	}
}

func TestHTTPSRouter(t *testing.T) {
	app := New()
	r := app.HTTPS()
	if h, _ := r.Lookup("GET", "/", nil); h != nil {
		t.Log("GET / should not return handler before registration")
		t.FailNow()
	}
	r.GET("/", 1)
	if h, _ := r.Lookup("GET", "/", nil); h == nil {
		t.Log("GET / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("GET", "/abc", nil); h != nil {
		t.Log("GET /abc should not return handler before registration")
		t.FailNow()
	}

	r.GET("/abc", "abc")

	if h, _ := r.Lookup("GET", "/abc", nil); h == nil {
		t.Log("GET /abc should return handler after registration")
		t.FailNow()
	}

	r.GET("/redir", app.ToTLSHandler())

	if h, _ := r.Lookup("GET", "/redir", nil); h == nil {
		t.Log("GET /abc should return handler after registration")
		t.FailNow()
	} else {
		defer func() {
			e := recover()
			if e != nil {
				t.Log("panic handled when testing /redir")
			}
		}()
		h(&Context{
			RequestCtx: &fasthttp.RequestCtx{},
		})
	}

	// POST

	if h, _ := r.Lookup("POST", "/", nil); h != nil {
		t.Log("POST / should not return handler before registration")
		t.FailNow()
	}
	r.POST("/", 1)
	if h, _ := r.Lookup("POST", "/", nil); h == nil {
		t.Log("POST / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("POST", "/abc", nil); h != nil {
		t.Log("POST /abc should not return handler before registration")
		t.FailNow()
	}

	r.POST("/abc", "abc")

	if h, _ := r.Lookup("POST", "/abc", nil); h == nil {
		t.Log("POST /abc should return handler after registration")
		t.FailNow()
	}

	// PUT
	if h, _ := r.Lookup("PUT", "/", nil); h != nil {
		t.Log("PUT / should not return handler before registration")
		t.FailNow()
	}
	r.PUT("/", 1)
	if h, _ := r.Lookup("PUT", "/", nil); h == nil {
		t.Log("PUT / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("PUT", "/abc", nil); h != nil {
		t.Log("PUT /abc should not return handler before registration")
		t.FailNow()
	}

	r.PUT("/abc", "abc")

	if h, _ := r.Lookup("PUT", "/abc", nil); h == nil {
		t.Log("PUT /abc should return handler after registration")
		t.FailNow()
	}

	// DELETE
	if h, _ := r.Lookup("DELETE", "/", nil); h != nil {
		t.Log("DELETE / should not return handler before registration")
		t.FailNow()
	}
	r.DELETE("/", 1)
	if h, _ := r.Lookup("DELETE", "/", nil); h == nil {
		t.Log("DELETE / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("DELETE", "/abc", nil); h != nil {
		t.Log("DELETE /abc should not return handler before registration")
		t.FailNow()
	}

	r.DELETE("/abc", "abc")

	if h, _ := r.Lookup("DELETE", "/abc", nil); h == nil {
		t.Log("DELETE /abc should return handler after registration")
		t.FailNow()
	}

	// HEAD
	if h, _ := r.Lookup("HEAD", "/", nil); h != nil {
		t.Log("HEAD / should not return handler before registration")
		t.FailNow()
	}
	r.HEAD("/", 1)
	if h, _ := r.Lookup("HEAD", "/", nil); h == nil {
		t.Log("HEAD / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("HEAD", "/abc", nil); h != nil {
		t.Log("HEAD /abc should not return handler before registration")
		t.FailNow()
	}

	r.HEAD("/abc", "abc")

	if h, _ := r.Lookup("HEAD", "/abc", nil); h == nil {
		t.Log("HEAD /abc should return handler after registration")
		t.FailNow()
	}

	// OPTIONS
	if h, _ := r.Lookup("OPTIONS", "/", nil); h != nil {
		t.Log("OPTIONS / should not return handler before registration")
		t.FailNow()
	}
	r.OPTIONS("/", 1)
	if h, _ := r.Lookup("OPTIONS", "/", nil); h == nil {
		t.Log("OPTIONS / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("OPTIONS", "/abc", nil); h != nil {
		t.Log("OPTIONS /abc should not return handler before registration")
		t.FailNow()
	}

	r.OPTIONS("/abc", "abc")

	if h, _ := r.Lookup("OPTIONS", "/abc", nil); h == nil {
		t.Log("OPTIONS /abc should return handler after registration")
		t.FailNow()
	}

	// PATCH
	if h, _ := r.Lookup("PATCH", "/", nil); h != nil {
		t.Log("PATCH / should not return handler before registration")
		t.FailNow()
	}
	r.PATCH("/", 1)
	if h, _ := r.Lookup("PATCH", "/", nil); h == nil {
		t.Log("PATCH / should return handler after registration")
		t.FailNow()
	}
	if h, _ := r.Lookup("PATCH", "/abc", nil); h != nil {
		t.Log("PATCH /abc should not return handler before registration")
		t.FailNow()
	}

	r.PATCH("/abc", "abc")

	if h, _ := r.Lookup("PATCH", "/abc", nil); h == nil {
		t.Log("PATCH /abc should return handler after registration")
		t.FailNow()
	}

	r.GET("/err", func(ctx *Context) error {
		return errors.New("test")
	})

	r.GET("/fasterr", func(ctx *fasthttp.RequestCtx) error {
		return errors.New("test")
	})

	if h, _ := r.Lookup("PATCH", "/abc", nil); h == nil {
		t.Log("PATCH /abc should return handler after registration")
		t.FailNow()
	}

	port := testutils.Port().NonRoot().Unused().Acquire()
	bindAddr := fmt.Sprintf(":%d", port)
	go func() {
		err := app.ListenAndServe(bindAddr)
		if err != nil {
			panic(err)
		}
	}()
	_, err := http.Get("http://127.0.0.1" + bindAddr) // just should not panic
	_ = err
	_, err = http.Get("http://127.0.0.1" + bindAddr) // just should not panic, twice
	_ = err
}
