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
	"github.com/valyala/fasthttp"
)

func (r *Router) getErrorHandler(h func(*Context) error) func(*Context) {
	return func(ctx *Context) {
		if err := h(ctx); err != nil {
			r.app.internalLog.WithField("url", ctx.URI()).Errorf("Error occurred: %s", err)
			ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		}
	}
}

func (r *Router) getGrameHandler(h func(*fasthttp.RequestCtx)) func(*Context) {
	return func(ctx *Context) {
		if ctx != nil {
			h(ctx.RequestCtx)
			return
		}
		h(new(fasthttp.RequestCtx))
	}
}

func (r *Router) getGrameDumbHandler(h func()) func(*Context) {
	return func(*Context) {
		h()
	}
}

func (r *Router) getGrameDumbErrorHandler(h func() error) func(*Context) {
	return func(ctx *Context) {
		if err := h(); err != nil {
			r.app.internalLog.WithField("url", ctx.URI()).Errorf("Error occurred: %s", err)
			ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		}
	}
}

func (r *Router) getGrameErrorHandler(h func(*fasthttp.RequestCtx) error) func(*Context) {
	return func(ctx *Context) {
		if err := h(ctx.RequestCtx); err != nil {
			r.app.internalLog.WithField("url", ctx.URI()).Errorf("Error occurred: %s", err)
			ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		}
	}
}

func (r *Router) initGrameCtx(ctx *fasthttp.RequestCtx) *Context {
	return &Context{
		Logger:     r.app.Logger,
		RequestCtx: ctx,
		App:        r.app,
	}
}

func (r *Router) initRouter() {
	if r.router == nil {
		r.router = newRouter()
	}
}

func (r *Router) getHandlerEncoder(h func() map[string]interface{}) func(*Context) {
	return func(ctx *Context) {
		r := h()
		if r == nil {
			ctx.SetStatusCode(fasthttp.StatusNoContent)
			return
		}
		if err := ctx.JSON(r); err != nil {
			ctx.jsonErrorLog(err)
		}
	}
}

func (r *Router) getCtxHandlerEncoder(h func(*Context) map[string]interface{}) func(*Context) {
	return func(ctx *Context) {
		r := h(ctx)
		if r == nil {
			ctx.SetStatusCode(fasthttp.StatusNoContent)
			return
		}
		if err := ctx.JSON(r); err != nil {
			ctx.jsonErrorLog(err)
		}
	}
}

func (r *Router) getHandlerEncoderErr(h func() (map[string]interface{}, error)) func(*Context) {
	return func(ctx *Context) {
		r, err := h()
		if err != nil {
			ctx.jsonErrorLog(err)
			return
		}
		if r == nil { // err == nil here
			ctx.SetStatusCode(fasthttp.StatusNoContent)
			return
		}
		if err = ctx.JSON(r); err != nil {
			ctx.jsonErrorLog(err)
		}
	}
}

func (r *Router) getCtxHandlerEncoderErr(h func(*Context) (map[string]interface{}, error)) func(*Context) {
	return func(ctx *Context) {
		r, err := h(ctx)
		if err != nil {
			ctx.jsonErrorLog(err)
			return
		}
		if r == nil { // err == nil here
			ctx.SetStatusCode(fasthttp.StatusNoContent)
			return
		}
		if err = ctx.JSON(r); err != nil {
			ctx.jsonErrorLog(err)
		}
	}
}

func (r *Router) getEfaceEncoder(h func() interface{}) func(*Context) {
	return func(ctx *Context) {
		r := h()
		if r == nil { // err == nil here
			ctx.SetStatusCode(fasthttp.StatusNoContent)
			return
		}
		if err := ctx.JSON(r); err != nil {
			ctx.jsonErrorLog(err)
		}
	}
}

func (r *Router) getEfaceErrEncoder(h func() (interface{}, error)) func(*Context) {
	return func(ctx *Context) {
		r, err := h()
		if err != nil {
			ctx.jsonErrorLog(err)
			return
		}
		if r == nil { // err == nil here
			ctx.SetStatusCode(fasthttp.StatusNoContent)
			return
		}
		if err = ctx.JSON(r); err != nil {
			ctx.jsonErrorLog(err)
		}
	}
}

func (r *Router) getEfaceCtxEncoder(h func(*Context) interface{}) func(*Context) {
	return func(ctx *Context) {
		r := h(ctx)
		if r == nil { // err == nil here
			ctx.SetStatusCode(fasthttp.StatusNoContent)
			return
		}
		if err := ctx.JSON(r); err != nil {
			ctx.jsonErrorLog(err)
		}
	}
}

func (r *Router) getEfaceCtxErrEncoder(h func(*Context) (interface{}, error)) func(*Context) {
	return func(ctx *Context) {
		r, err := h(ctx)
		if err != nil {
			ctx.jsonErrorLog(err)
			return
		}
		if r == nil { // err == nil here
			ctx.SetStatusCode(fasthttp.StatusNoContent)
			return
		}
		if err = ctx.JSON(r); err != nil {
			ctx.jsonErrorLog(err)
		}
	}
}
