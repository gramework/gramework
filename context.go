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

// @TODO: add more
var ctypes = []string{
	jsonCT,
	xmlCT,
}

// Encode automatically determies accepted formats
// and choose preferred one
func (c *Context) Encode(v interface{}) (sentType string, err error) {
	accept := c.Request.Header.Peek(acceptHeader)
	accepted := acceptParser.Parse(BytesToString(accept))

	sentType, err = accepted.Negotiate(ctypes...)
	if err != nil {
		return
	}

	switch sentType {
	case jsonCT:
		c.JSON(v)
	case xmlCT:
		c.XML(v)
	}

	return
}

// XML sends text/xml content type (see rfc3023, sec 3) and xml-encoded value to client
func (c *Context) XML(v interface{}) error {
	c.SetContentType(xmlCT)
	b, err := c.ToXML(v)
	if err != nil {
		return err
	}

	c.Write(b)
	return nil
}

// ToXML encodes xml-encoded value to client
func (c *Context) ToXML(v interface{}) ([]byte, error) {
	b := bytes.NewBuffer(nil)
	err := xml.NewEncoder(b).Encode(v)
	return b.Bytes(), err
}

// GETKeys returns GET parameters keys
func (c *Context) GETKeys() []string {
	res := []string{}
	c.Request.URI().QueryArgs().VisitAll(func(key, value []byte) {
		res = append(res, string(key))
	})
	return res
}

// GETKeysBytes returns GET parameters keys as []byte
func (c *Context) GETKeysBytes() [][]byte {
	res := [][]byte{}
	c.Request.URI().QueryArgs().VisitAll(func(key, value []byte) {
		res = append(res, key)
	})
	return res
}

// GETParams returns GET parameters
func (c *Context) GETParams() map[string][]string {
	res := map[string][]string{}
	c.Request.URI().QueryArgs().VisitAll(func(key, value []byte) {
		res[string(key)] = append(res[string(key)], string(value))
	})
	return res
}

// GETParam returns GET parameter by name
func (c *Context) GETParam(argName string) []string {
	res := c.GETParams()
	if param, ok := res[argName]; ok {
		return param
	}
	return []string{}
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

// ToTLS redirects user to HTTPS scheme
func (c *Context) ToTLS() {
	u := c.URI()
	u.SetScheme(https)
	c.Redirect(u.String(), redirectCode)
}

const (
	redirectCode                      = 301
	temporaryRedirectCode             = 307
	zero                              = 0
	one                               = 1
	https                             = "https"
	corsAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	corsAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	corsAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	corsAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	methods                           = "GET,PUT,POST,DELETE"
	corsCType                         = "Content-Type, *"
	trueStr                           = "true"
	jsonCT                            = "application/json;charset=utf8"
	hOrigin                           = "Origin"
	forbidden                         = "Forbidden"
	forbiddenCode                     = 403
)

// CORS enables CORS in the current context
func (c *Context) CORS(domains ...string) *Context {
	origins := make([]string, 0)
	if len(domains) > 0 {
		origins = domains
	} else if headerOrigin := c.Request.Header.Peek(hOrigin); headerOrigin != nil && len(headerOrigin) > 0 {
		origins = append(origins, string(headerOrigin))
	} else {
		origins = append(origins, string(c.Request.URI().Host()))
	}

	c.Response.Header.Set(corsAccessControlAllowOrigin, strings.Join(origins, " "))
	c.Response.Header.Set(corsAccessControlAllowMethods, methods)
	c.Response.Header.Set(corsAccessControlAllowHeaders, corsCType)
	c.Response.Header.Set(corsAccessControlAllowCredentials, trueStr)

	return c
}

// Forbidden send 403 Forbidden error
func (c *Context) Forbidden() {
	c.Error(forbidden, forbiddenCode)
}

// JSON serializes and writes a json-formatted response to user
func (c *Context) JSON(v interface{}) error {
	c.SetContentType(jsonCT)
	b, err := c.ToJSON(v)
	c.Write(b)
	return err
}

// ToJSON serializes and returns the result
func (c *Context) ToJSON(v interface{}) ([]byte, error) {
	b := bytes.NewBuffer(nil)
	err := json.NewEncoder(b).Encode(v)
	return b.Bytes(), err
}

// UnJSONBytes serializes and writes a json-formatted response to user
func (c *Context) UnJSONBytes(b []byte, v ...interface{}) (interface{}, error) {
	return UnJSONBytes(b, v...)
}

// UnJSON deserializes JSON request body to given variable pointer
func (c *Context) UnJSON(v interface{}) error {
	return json.NewDecoder(bytes.NewReader(c.Request.Body())).Decode(&v)
}

// UnJSONBytes serializes and writes a json-formatted response to user
func UnJSONBytes(b []byte, v ...interface{}) (interface{}, error) {
	if len(v) == 0 {
		var res interface{}
		err := json.NewDecoder(bytes.NewReader(b)).Decode(&res)
		return res, err
	}
	err := json.NewDecoder(bytes.NewReader(b)).Decode(&v[0])
	return v[0], err
}

// Err500 sets Internal Server Error status
func (c *Context) Err500(message ...interface{}) *Context {
	c.SetStatusCode(fasthttp.StatusInternalServerError)
	for k := range message {
		switch v := message[k].(type) {
		case string:
			c.WriteString(v)
		case error:
			c.Writef("%s", v)
		default:
			c.Writef("%v", v)
		}
	}
	return c
}

// JSONError sets Internal Server Error status,
// serializes and writes a json-formatted response to user
func (c *Context) JSONError(v interface{}) error {
	c.Err500()
	return c.JSON(v)
}
