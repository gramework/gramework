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
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/gocarina/gocsv"
	"github.com/gramework/runtimer"
)

// @TODO: add more
var ctypes = []string{
	jsonCT,
	xmlCT,
	csvCT,
}

// ContextFromValue returns gramework.Context from context.Context value from gramework.ContextKey
// in a more effective way, than standard eface.(*SomeType).
// WARNING: this function may return nil, if ctx has no gramework.Context stored or ctx is nil.
// This function will give a warning if you call it with nil context.Context.
func ContextFromValue(ctx context.Context) *Context {
	if ctx == nil {
		internalLog.Warn("ContextFromValue was called with nil context.Context, returning nil")
		return nil
	}
	return (*Context)(runtimer.GetEfaceDataPtr(ctx.Value(ContextKey)))
}

// MWKill kills current context and stop any user-defined processing.
// This function intented for use in middlewares.
func (ctx *Context) MWKill() {
	ctx.middlewareKilledReq = true
}

// SubPrefixes returns list of router's prefixes that was created using .Sub() feature
func (ctx *Context) SubPrefixes() []string {
	return ctx.subPrefixes
}

// ContentType returns Content-Type header for current request
func (ctx *Context) ContentType() string {
	return string(ctx.Request.Header.Peek(contentType))
}

// ToContext returns context.Context with gramework.Context stored
// in context values as a pointer (see gramework.ContextKey to receive and use this value).
//
// By default this func will extend context.Background(), if parentCtx is not provided.
func (ctx *Context) ToContext(parentCtx ...context.Context) context.Context {
	if len(parentCtx) > 0 {
		return context.WithValue(parentCtx[0], ContextKey, ctx)
	}

	return context.WithValue(context.Background(), ContextKey, ctx)
}

// RouteArg returns an argument value as a string or empty string
func (ctx *Context) RouteArg(argName string) string {
	v, err := ctx.RouteArgErr(argName)
	if err != nil {
		return emptyString
	}
	return v
}

// ToCSV encodes csv-encoded value to client
func (ctx *Context) ToCSV(v interface{}) ([]byte, error) {
	return gocsv.MarshalBytes(v)
}

// ToXML encodes xml-encoded value to client
func (ctx *Context) ToXML(v interface{}) ([]byte, error) {
	b := bytes.NewBuffer(nil)
	err := xml.NewEncoder(b).Encode(v)
	return b.Bytes(), err
}

// GETKeys returns GET parameters keys (query args)
func (ctx *Context) GETKeys() []string {
	var res []string
	ctx.Request.URI().QueryArgs().VisitAll(func(key, value []byte) {
		res = append(res, string(key))
	})
	return res
}

// GETKeysBytes returns GET parameters keys (query args) as []byte
func (ctx *Context) GETKeysBytes() [][]byte {
	var res [][]byte
	ctx.Request.URI().QueryArgs().VisitAll(func(key, value []byte) {
		res = append(res, key)
	})
	return res
}

// GETParams returns GET parameters (query args)
func (ctx *Context) GETParams() map[string][]string {
	res := make(map[string][]string)
	ctx.Request.URI().QueryArgs().VisitAll(func(key, value []byte) {
		res[string(key)] = append(res[string(key)], string(value))
	})
	return res
}

// GETParam returns GET parameter (query arg) by name
func (ctx *Context) GETParam(argName string) []string {
	res := ctx.GETParams()
	if param, ok := res[argName]; ok {
		return param
	}
	return nil
}

// RouteArgErr returns an argument value as a string or empty string
// and ErrArgNotFound if argument was not found
func (ctx *Context) RouteArgErr(argName string) (string, error) {
	i := ctx.UserValue(argName)
	if i == nil {
		return emptyString, ErrArgNotFound
	}
	switch value := i.(type) {
	case string:
		return value, nil
	default:
		return fmt.Sprintf(fmtV, i), nil
	}
}

// ToTLS redirects user to HTTPS scheme
func (ctx *Context) ToTLS() {
	u := ctx.URI()
	u.SetScheme(https)
	ctx.Redirect(u.String(), redirectCode)
}

// Forbidden send 403 Forbidden error
func (ctx *Context) Forbidden() {
	ctx.Error(forbidden, forbiddenCode)
}

// ToJSON serializes v and returns the result
func (ctx *Context) ToJSON(v interface{}) ([]byte, error) {
	b := bytes.NewBuffer(nil)
	err := json.NewEncoder(b).Encode(v)
	return b.Bytes(), err
}

// UnJSONBytes deserializes JSON request body to given variable pointer or allocates a new one.
// Returns resulting data and error. One of them may be nil.
func (ctx *Context) UnJSONBytes(b []byte, v ...interface{}) (interface{}, error) {
	return UnJSONBytes(b, v...)
}

// UnJSON deserializes JSON request body to given variable pointer
func (ctx *Context) UnJSON(v interface{}) error {
	return json.NewDecoder(bytes.NewReader(ctx.Request.Body())).Decode(&v)
}

// UnJSONBytes deserializes JSON request body to given variable pointer or allocates a new one.
// Returns resulting data and error. One of them may be nil.
func UnJSONBytes(b []byte, v ...interface{}) (interface{}, error) {
	if len(v) == 0 {
		var res interface{}
		err := json.NewDecoder(bytes.NewReader(b)).Decode(&res)
		return res, err
	}
	err := json.NewDecoder(bytes.NewReader(b)).Decode(&v[0])
	return v[0], err
}

func (ctx *Context) jsonErrorLog(v interface{}) {
	ctx.Err500()
	if err := ctx.JSON(v); err != nil {
		ctx.Logger.WithError(err).Error("JSONError err")
	}
}

// RequestID return request ID for current context's request
func (ctx *Context) RequestID() string {
	return ctx.requestID
}
