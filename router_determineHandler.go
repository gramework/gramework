// Copyright 2017 Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package gramework

import "github.com/valyala/fasthttp"

func (r *Router) determineHandler(handler interface{}) func(*Context) {
	switch h := handler.(type) {
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
		r.app.Logger.Warnf("Unknown handler type: %T, serving fmt.Sprintf(%%v)", h)
		return r.getFmtVHandler(h)
	}
}
