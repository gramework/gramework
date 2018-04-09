package gramework

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	acceptParser "github.com/kirillDanshin/go-accept-headers"
	"github.com/valyala/fasthttp"
)

// @TODO: add more
var ctypes = []string{
	jsonCT,
	xmlCT,
}

// Writef is a fmt.Fprintf(context, format, a...) shortcut
func (ctx *Context) Writef(format string, a ...interface{}) {
	fmt.Fprintf(ctx, format, a...)
}

// Writeln is a fmt.Fprintln(context, format, a...) shortcut
func (ctx *Context) Writeln(a ...interface{}) {
	fmt.Fprintln(ctx, a...)
}

// RouteArg returns an argument value as a string or empty string
func (ctx *Context) RouteArg(argName string) string {
	v, err := ctx.RouteArgErr(argName)
	if err != nil {
		return emptyString
	}
	return v
}

// Encode automatically determies accepted formats
// and choose preferred one
func (ctx *Context) Encode(v interface{}) (sentType string, err error) {
	accept := ctx.Request.Header.Peek(acceptHeader)
	accepted := acceptParser.Parse(BytesToString(accept))

	sentType, err = accepted.Negotiate(ctypes...)
	if err != nil {
		return
	}

	switch sentType {
	case jsonCT:
		err = ctx.JSON(v)
	case xmlCT:
		err = ctx.XML(v)
	}

	return
}

// XML sends text/xml content type (see rfc3023, sec 3) and xml-encoded value to client
func (ctx *Context) XML(v interface{}) error {
	ctx.SetContentType(xmlCT)
	b, err := ctx.ToXML(v)
	if err != nil {
		return err
	}

	_, err = ctx.Write(b)
	return err
}

// ToXML encodes xml-encoded value to client
func (ctx *Context) ToXML(v interface{}) ([]byte, error) {
	b := bytes.NewBuffer(nil)
	err := xml.NewEncoder(b).Encode(v)
	return b.Bytes(), err
}

// GETKeys returns GET parameters keys
func (ctx *Context) GETKeys() []string {
	var res []string
	ctx.Request.URI().QueryArgs().VisitAll(func(key, value []byte) {
		res = append(res, string(key))
	})
	return res
}

// GETKeysBytes returns GET parameters keys as []byte
func (ctx *Context) GETKeysBytes() [][]byte {
	var res [][]byte
	ctx.Request.URI().QueryArgs().VisitAll(func(key, value []byte) {
		res = append(res, key)
	})
	return res
}

// GETParams returns GET parameters
func (ctx *Context) GETParams() map[string][]string {
	res := make(map[string][]string)
	ctx.Request.URI().QueryArgs().VisitAll(func(key, value []byte) {
		res[string(key)] = append(res[string(key)], string(value))
	})
	return res
}

// GETParam returns GET parameter by name
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

// HTML sets HTML content type
func (ctx *Context) HTML() *Context {
	ctx.SetContentType(htmlCT)
	return ctx
}

// ToTLS redirects user to HTTPS scheme
func (ctx *Context) ToTLS() {
	u := ctx.URI()
	u.SetScheme(https)
	ctx.Redirect(u.String(), redirectCode)
}

// CORS enables CORS in the current context
func (ctx *Context) CORS(domains ...string) *Context {
	origins := make([]string, 0)
	if len(domains) > 0 {
		origins = domains
	} else if headerOrigin := ctx.Request.Header.Peek(hOrigin); len(headerOrigin) > 0 {
		origins = append(origins, string(headerOrigin))
	} else {
		origins = append(origins, string(ctx.Request.URI().Host()))
	}

	ctx.Response.Header.Set(corsAccessControlAllowOrigin, strings.Join(origins, " "))
	ctx.Response.Header.Set(corsAccessControlAllowMethods, methods)
	ctx.Response.Header.Set(corsAccessControlAllowHeaders, corsCType)
	ctx.Response.Header.Set(corsAccessControlAllowCredentials, trueStr)

	return ctx
}

// Forbidden send 403 Forbidden error
func (ctx *Context) Forbidden() {
	ctx.Error(forbidden, forbiddenCode)
}

// JSON serializes and writes a json-formatted response to user
func (ctx *Context) JSON(v interface{}) error {
	ctx.SetContentType(jsonCT)
	b, err := ctx.ToJSON(v)
	if err != nil {
		return err
	}

	_, err = ctx.Write(b)
	return err
}

// ToJSON serializes and returns the result
func (ctx *Context) ToJSON(v interface{}) ([]byte, error) {
	b := bytes.NewBuffer(nil)
	err := json.NewEncoder(b).Encode(v)
	return b.Bytes(), err
}

// UnJSONBytes serializes and writes a json-formatted response to user
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

// BadRequest sends HTTP/1.1 400 Bad Request
func (ctx *Context) BadRequest(err ...error) {
	ctx.Error(badRequest, fasthttp.StatusBadRequest)
	if len(err) == 1 {
		ctx.Writeln(err[0])
	}
}

// Err500 sets Internal Server Error status
func (ctx *Context) Err500(message ...interface{}) *Context {
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	for k := range message {
		switch v := message[k].(type) {
		case string:
			ctx.WriteString(v)
		case error:
			ctx.Writef("%s", v)
		default:
			ctx.Writef("%v", v)
		}
	}
	return ctx
}

// JSONError sets Internal Server Error status,
// serializes and writes a json-formatted response to user
func (ctx *Context) JSONError(v interface{}) error {
	ctx.Err500()
	return ctx.JSON(v)
}
