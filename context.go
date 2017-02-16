package gramework

import (
	"encoding/json"
	"fmt"

	"bytes"

	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
)

var buffer bytebufferpool.Pool

// Writef is a fmt.Fprintf(context, format, a...) shortcut
func (c *Context) Writef(format string, a ...interface{}) {
	fmt.Fprintf(c, format, a...)
}

// Writeln is a fmt.Fprintln(context, format, a...) shortcut
func (c *Context) Writeln(a ...interface{}) {
	fmt.Fprintln(c, a...)
}

// RouteArg returns an argument value as a string or empty string
func (c *Context) RouteArg(argName string) string {
	v, err := c.RouteArgErr(argName)
	if err != nil {
		return emptyString
	}
	return v
}

// RouteArgErr returns an argument value as a string or empty string
// and ErrArgNotFound if argument was not found
func (c *Context) RouteArgErr(argName string) (string, error) {
	i := c.UserValue(argName)
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

// HTML sets HTML content type
func (c *Context) HTML() *Context {
	c.SetContentType(htmlCT)
	return c
}

const (
	corsAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	corsAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	corsAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	corsAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	methods                           = "GET,PUT,POST,DELETE"
	corsCType                         = "Content-Type, *"
	trueStr                           = "true"
	jsonCT                            = "application/json"
)

// CORS enables CORS in the current context
func (c *Context) CORS() *Context {
	c.Response.Header.Set(corsAccessControlAllowOrigin, string(c.URI().Host()))
	c.Response.Header.Set(corsAccessControlAllowMethods, methods)
	c.Response.Header.Set(corsAccessControlAllowHeaders, corsCType)
	c.Response.Header.Set(corsAccessControlAllowCredentials, trueStr)

	return c
}

// JSON serializes and writes a json-formatted response to user
func (c *Context) JSON(v interface{}) error {
	c.SetContentType(jsonCT)
	w := buffer.Get()
	err := json.NewEncoder(w).Encode(v)
	c.Write(w.Bytes())
	buffer.Put(w)
	return err
}

// UnJSONBytes serializes and writes a json-formatted response to user
func (c *Context) UnJSONBytes(b []byte, v interface{}) (interface{}, error) {
	var res interface{}
	err := json.NewDecoder(bytes.NewReader(b)).Decode(&res)
	return res, err
}

// Err500 sets Internal Server Error status
func (c *Context) Err500() *Context {
	c.SetStatusCode(fasthttp.StatusInternalServerError)

	return c
}

// JSONError sets Internal Server Error status,
// serializes and writes a json-formatted response to user
func (c *Context) JSONError(v interface{}) error {
	c.Err500()
	return c.JSON(v)
}
