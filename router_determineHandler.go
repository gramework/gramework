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
	"errors"
	"reflect"
	"strings"

	"github.com/valyala/fasthttp"
)

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
	case HTML:
		return r.getHTMLServer(h)
	case JSON:
		return r.getJSONServer(h)
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
		rv := reflect.ValueOf(h)
		if rv.Kind() == reflect.Func {
			handler, err := r.getCachedReflectHandler(h)
			if err != nil {
				r.app.internalLog.WithError(err).Fatal("Unsupported reflect handler signature")
			}

			return handler
		}
		r.app.internalLog.Warnf("Unknown handler type: %T, serving fmt.Sprintf(%%v)", h)
		return r.getFmtVHandler(h)
	}
}

type reflectDecodedBodyRecv struct {
	idx int
	t   reflect.Type
}

func (r *Router) getCachedReflectHandler(h interface{}) (func(*Context), error) {
	funcT := reflect.TypeOf(h)
	if funcT.IsVariadic() {
		return nil, errors.New("could not process variadic reflect handler")
	}

	results := funcT.NumOut()
	if results > 2 {
		return nil, errors.New("reflect handler output should be one of (any), (any, error), (error) or ()")
	}

	params := funcT.NumIn()
	decodedBodyRecv := []reflectDecodedBodyRecv{}
	ctxRecv := -1

	checkForErrorAt := -1
	encodeDataAt := -1

	for i := 0; i < params; i++ {
		p := funcT.In(i)
		if strings.Contains(p.String(), "*gramework.Context") {
			ctxRecv = i
			continue
		}
		decodedBodyRecv = append(decodedBodyRecv, reflectDecodedBodyRecv{
			idx: i,
			t:   p,
		})
	}

	for i := 0; i < results; i++ {
		r := funcT.Out(i)
		println(r.String())

		if r.String() == "error" {
			if i == 0 && results > 1 {
				return nil, errors.New("reflect handler output should be one of (any), (any, error), (error) or ()")
			}

			checkForErrorAt = i
			continue
		}

		if encodeDataAt >= 0 {
			return nil, errors.New("reflect handler output should be one of (any), (any, error), (error) or ()")
		}

		encodeDataAt = i
	}

	funcV := reflect.ValueOf(h)

	handler := func(ctx *Context) {
		callParams := make([]reflect.Value, params)
		if len(decodedBodyRecv) > 0 {
			unsupportedBodyType := true
			for i := range decodedBodyRecv {
				decoded := reflect.New(decodedBodyRecv[i].t).Interface()
				if jsonErr := ctx.UnJSON(decoded); jsonErr == nil {
					unsupportedBodyType = false
					decodedV := reflect.ValueOf(decoded)

					callParams[decodedBodyRecv[i].idx] = decodedV.Elem()
				} else {
					callParams[decodedBodyRecv[i].idx] = reflect.Zero(decodedBodyRecv[i].t)
				}
			}

			if unsupportedBodyType {
				ctx.SetStatusCode(500)
				ctx.Logger.Error("unsupported body type")
				return
			}
		}
		if ctxRecv >= 0 {
			callParams[ctxRecv] = reflect.ValueOf(ctx)
		}

		res := funcV.Call(callParams)
		shouldProcessErr := false
		shouldProcessReturn := false
		var err error
		if checkForErrorAt >= 0 && !res[checkForErrorAt].IsNil() {
			resErr, ok := res[checkForErrorAt].Interface().(error)
			if ok {
				err = resErr
			} else {
				err = errUnknown
			}
			shouldProcessErr = true
		}

		var v interface{}
		if encodeDataAt >= 0 {
			v = res[encodeDataAt].Interface()
			shouldProcessReturn = true
		}
		if shouldProcessErr {
			if err != nil {
				ctx.jsonErrorLog(err)
				return
			}
		}
		if shouldProcessReturn {
			if v == nil { // err == nil here
				ctx.SetStatusCode(fasthttp.StatusNoContent)
				return
			}
			if err = ctx.JSON(v); err != nil {
				ctx.jsonErrorLog(err)
			}
		}
		return
	}

	return handler, nil
}

var errUnknown = errors.New("Unknown Server Error")
