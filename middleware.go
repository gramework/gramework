// Copyright 2017 Kirill Danshin and Gramework contributors
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

	"github.com/valyala/fasthttp"
)

var (
	// ErrEmptyMiddleware can be returned by App.Use*, if middleware is nil
	ErrEmptyMiddleware = errors.New("can't use nil middleware")

	// ErrUnsupportedMiddlewareType can be returned by App.Use*, if middleware type is unsupported
	ErrUnsupportedMiddlewareType = errors.New("unsupported middleware type")
)

// CORSMiddleware provides gramework handler with ctx.CORS() call
func (app *App) CORSMiddleware(domains ...string) func(*Context) {
	return func(ctx *Context) {
		ctx.CORS(domains...)
	}
}

// Use the middleware before request processing
func (app *App) Use(middleware interface{}) error {
	if middleware == nil {
		return ErrEmptyMiddleware
	}
	processor, err := app.middlewareProcessor(middleware)
	app.middlewaresMu.Lock()
	if err == nil {
		app.middlewares = append(app.middlewares, processor)
	}
	app.middlewaresMu.Unlock()

	return err
}

// UsePre registers middleware before any other middleware. Use only for metrics or access control!
func (app *App) UsePre(middleware interface{}) error {
	if middleware == nil {
		return ErrEmptyMiddleware
	}
	processor, err := app.middlewareProcessor(middleware)
	app.preMiddlewaresMu.Lock()
	if err == nil {
		app.preMiddlewares = append(app.preMiddlewares, processor)
	}
	app.preMiddlewaresMu.Unlock()

	return err
}

// UseAfterRequest the middleware after request processing
func (app *App) UseAfterRequest(middleware interface{}) error {
	if middleware == nil {
		return ErrEmptyMiddleware
	}

	processor, err := app.middlewareProcessor(middleware)
	app.middlewaresAfterRequestMu.Lock()
	if err == nil {
		app.middlewaresAfterRequest = append(app.middlewaresAfterRequest, processor)
	}
	app.middlewaresAfterRequestMu.Unlock()

	return nil
}

func (app *App) middlewareProcessor(middleware interface{}) (func(*Context), error) {
	// we can register middlewares slowly to serve requests faster
	switch m := middleware.(type) {
	case func():
		return func(*Context) {
			m()
		}, nil
	case func(*Context):
		// if middleware is that type, we can just return
		// the middleware itself, to save some resources
		// required to run the function via closures
		return m, nil
	case func(*Context) error:
		return func(ctx *Context) {
			if err := m(ctx); err != nil {
				// if error occurred, we can stop processing even slowly
				ctx.Logger.Errorf("Middleware error: %s", err)
				ctx.Err500()
				ctx.middlewaresShouldStopProcessing = true
			}
		}, nil
	case func(*fasthttp.RequestCtx):
		return func(ctx *Context) {
			m(ctx.RequestCtx)
		}, nil
	}
	return nil, ErrUnsupportedMiddlewareType
}
