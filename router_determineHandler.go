// Copyright 2017-present Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package gramework

import "github.com/valyala/fasthttp"

type reqHandlerDefault interface {
	Handler(*Context)
}

type reqHandlerWithError interface {
	Handler(*Context) error
}

type reqHandlerWithEfaceError interface {
	Handler(*Context) (interface{}, error)
}

type reqHandlerWithEface interface {
	Handler(*Context) interface{}
}

type reqHandlerNoCtx interface {
	Handler()
}

type reqHandlerWithErrorNoCtx interface {
	Handler() error
}

type reqHandlerWithEfaceErrorNoCtx interface {
	Handler() (interface{}, error)
}

type reqHandlerWithEfaceNoCtx interface {
	Handler() interface{}
}

func (r *Router) determineHandler(handler interface{}) func(*Context) {
	// copy handler, we don't want to mutate our arguments
	rawHandler := handler

	// prepare handler in case if it one of our supported interfaces
	switch h := handler.(type) {
	case reqHandlerDefault:
		rawHandler = h.Handler
	case reqHandlerWithError:
		rawHandler = h.Handler
	case reqHandlerWithEfaceError:
		rawHandler = h.Handler
	case reqHandlerWithEface:
		rawHandler = h.Handler
	case reqHandlerNoCtx:
		rawHandler = h.Handler
	case reqHandlerWithErrorNoCtx:
		rawHandler = h.Handler
	case reqHandlerWithEfaceErrorNoCtx:
		rawHandler = h.Handler
	case reqHandlerWithEfaceNoCtx:
		rawHandler = h.Handler
	}

	// finally, process the handler
	switch h := rawHandler.(type) {
	case func(*Context):
		return h
	case RequestHandler:
		return h
	case func(*Context) error:
		return r.getErrorHandler(h)
	case func(*fasthttp.RequestCtx):
		return r.getGrameHandler(h)
	case func(*fasthttp.RequestCtx) error:
		return r.getGrameErrorHandler(h)
	case func() interface{}:
		return r.getEfaceEncoder(h)
	case func() (interface{}, error):
		return r.getEfaceErrEncoder(h)
	case func(*Context) interface{}:
		return r.getEfaceCtxEncoder(h)
	case func(*Context) (interface{}, error):
		return r.getEfaceCtxErrEncoder(h)
	case string:
		return r.getStringServer(h)
	case []byte:
		return r.getBytesServer(h)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return r.getFmtDHandler(h)
	case float32, float64:
		return r.getFmtFHandler(h)
	case func():
		return r.getGrameDumbHandler(h)
	case func() error:
		return r.getGrameDumbErrorHandler(h)
	case func() string:
		return r.getEFuncStrHandler(h)
	case func() map[string]interface{}:
		return r.getHandlerEncoder(h)
	case func(*Context) map[string]interface{}:
		return r.getCtxHandlerEncoder(h)
	case func() (map[string]interface{}, error):
		return r.getHandlerEncoderErr(h)
	case func(*Context) (map[string]interface{}, error):
		return r.getCtxHandlerEncoderErr(h)
	default:
		r.app.internalLog.Warnf("Unknown handler type: %T, serving fmt.Sprintf(%%v)", h)
		return r.getFmtVHandler(h)
	}
}
