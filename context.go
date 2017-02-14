package gramework

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/bytebufferpool"
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
	w := buffer.Get()
	err := json.NewEncoder(w).Encode(v)
	c.Write(w.Bytes())
	buffer.Put(w)
	return err
}
