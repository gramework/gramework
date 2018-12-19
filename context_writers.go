package gramework

import (
	"fmt"
	"strings"

	acceptParser "github.com/kirillDanshin/go-accept-headers"
	"github.com/valyala/fasthttp"
)

// Encode automatically determines accepted formats
// and choose preferred one
func (ctx *Context) Encode(v interface{}) (string, error) {
	accept := ctx.Request.Header.Peek(acceptHeader)
	accepted := acceptParser.Parse(BytesToString(accept))

	sentType, err := accepted.Negotiate(ctypes...)
	if err != nil {
		return emptyString, err
	}

	switch sentType {
	case jsonCT:
		err = ctx.JSON(v)
	case xmlCT:
		err = ctx.XML(v)
	case csvCT:
		err = ctx.CSV(v)
	}

	return sentType, err
}

// Writef is a fmt.Fprintf(context, format, a...) shortcut
func (ctx *Context) Writef(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(ctx, format, a...)
}

// CSV sends text/csv content type (see rfc4180, sec 3) and csv-encoded value to client
func (ctx *Context) CSV(v interface{}) error {
	ctx.SetContentType(csvCT)

	b, err := ctx.ToCSV(v)
	if err != nil {
		return err
	}
	_, err = ctx.Write(b)
	return err
}

// Writeln is a fmt.Fprintln(context, format, a...) shortcut
func (ctx *Context) Writeln(a ...interface{}) (int, error) {
	return fmt.Fprintln(ctx, a...)
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

// HTML sets HTML content type
func (ctx *Context) HTML(src ...string) *Context {
	ctx.SetContentType(htmlCT)
	if len(src) > 0 {
		ctx.WriteString(src[0])
	}
	return ctx
}

// CORS enables CORS in the current context
func (ctx *Context) CORS(domains ...string) *Context {
	var origins []string
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

// BadRequest sends HTTP/1.1 400 Bad Request
func (ctx *Context) BadRequest(err ...error) {
	e := badRequest
	if len(err) > 0 {
		e = err[0].Error()
	}

	ctx.Error(e, fasthttp.StatusBadRequest)
}

// Err500 sets Internal Server Error status
func (ctx *Context) Err500(message ...interface{}) *Context {
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	for k := range message {
		switch v := message[k].(type) {
		case string:
			_, err := ctx.WriteString(v)
			if err != nil {
				ctx.Logger.WithError(err).Error("Err500 serving error")
			}
		case error:
			ctx.Writef(fmtS, v)
		default:
			ctx.Writef(fmtV, v)
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
