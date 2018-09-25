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
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/apex/log"
	"github.com/valyala/fasthttp"
)

// JSON register internal handler that sets json content type
// and serves given handler with GET method
func (r *Router) JSON(route string, handler interface{}) *Router {
	h := r.determineHandler(handler)
	r.GET(route, jsonHandler(h))

	return r
}

// GET registers a handler for a GET request to the given route
func (r *Router) GET(route string, handler interface{}) *Router {
	r.Handle(MethodGET, route, handler)
	return r
}

// Forbidden serves 403 on route it registered on
func (r *Router) Forbidden(ctx *Context) {
	ctx.Forbidden()
}

// DELETE registers a handler for a DELETE request to the given route
func (r *Router) DELETE(route string, handler interface{}) *Router {
	r.Handle(MethodDELETE, route, handler)
	return r
}

// HEAD registers a handler for a HEAD request to the given route
func (r *Router) HEAD(route string, handler interface{}) *Router {
	r.Handle(MethodHEAD, route, handler)
	return r
}

// OPTIONS registers a handler for a OPTIONS request to the given route
func (r *Router) OPTIONS(route string, handler interface{}) *Router {
	r.Handle(MethodOPTIONS, route, handler)
	return r
}

// PUT registers a handler for a PUT request to the given route
func (r *Router) PUT(route string, handler interface{}) *Router {
	r.Handle(MethodPUT, route, handler)
	return r
}

// POST registers a handler for a POST request to the given route
func (r *Router) POST(route string, handler interface{}) *Router {
	r.Handle(MethodPOST, route, handler)
	return r
}

// PATCH registers a handler for a PATCH request to the given route
func (r *Router) PATCH(route string, handler interface{}) *Router {
	r.Handle(MethodPATCH, route, handler)
	return r
}

// ServeFile serves a file on a given route
func (r *Router) ServeFile(route, file string) *Router {
	r.Handle(MethodGET, route, func(ctx *Context) {
		ctx.SendFile(file)
	})
	return r
}

// SPAIndex serves an index file or handler on any unregistered route
func (r *Router) SPAIndex(pathOrHandler interface{}) *Router {
	switch v := pathOrHandler.(type) {
	case string:
		r.NotFound(func(ctx *Context) {
			ctx.HTML()
			ctx.SendFile(v)
		})
	default:
		r.NotFound(r.determineHandler(v))
	}
	return r
}

// Sub let you quickly register subroutes with given prefix
// like app.Sub("v1").GET("route", "hi"), that give you /v1/route
// registered
func (r *Router) Sub(path string) *SubRouter {
	return &SubRouter{
		prefix:   path,
		parent:   r,
		prefixes: []string{path},
	}
}

func (r *Router) handleReg(method, route string, handler interface{}, prefixes []string) {
	r.initRouter()
	r.app.internalLog.Debugf("registering %s %s", method, route)
	typedHandler := r.determineHandler(handler)
	for prefix := range r.app.protectedPrefixes {
		if strings.HasPrefix(strings.TrimLeft(route, "/"), strings.TrimLeft(prefix, "/")) {
			r.app.internalLog.
				WithField("route", route).
				WithField("method", method).
				Info("[Gramework Protection] Protection enabled for a new route")
			r.app.protectedEndpoints[route] = struct{}{}
			typedHandler = r.app.protectionMiddleware(typedHandler)
			break
		}
	}
	r.router.Handle(method, route, typedHandler, prefixes)
}

func (r *Router) getEFuncStrHandler(h func() string) func(*Context) {
	return func(ctx *Context) {
		ctx.WriteString(h())
	}
}

func handlerName(h interface{}) string {
	v := reflect.ValueOf(h)
	if v.Kind() != reflect.Func {
		return fmt.Sprintf("<raw %T>", h)
	}
	funcDesc := runtime.FuncForPC(v.Pointer())
	file, line := funcDesc.FileLine(v.Pointer())
	pathidx := strings.Index(file, "/go/src/")
	if strings.Contains(file, "/go/src") && len(file) > pathidx+len("/go/src/") {
		file = file[strings.Index(file, "/go/src/")+8:]
	}
	name := fmt.Sprintf("%s@%s:%v", funcDesc.Name(), file, line)
	return name
}

// Handle registers a new request handle with the given path and method.
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut functions can be used.
// This function is intended for bulk loading and to allow the usage of less frequently used,
// non-standardized or custom methods (e.g. for internal communication with a proxy).
func (r *Router) Handle(method, route string, handler interface{}) *Router {
	r.app.internalLog.WithFields(log.Fields{
		"handler": handlerName(handler),
		"method":  method,
		"route":   route,
	}).Debug("registering route")
	r.handleReg(method, route, handler, nil)
	return r
}

func (r *Router) getFmtVHandler(v interface{}) func(*Context) {
	cache := []byte(fmt.Sprintf("%v", v))
	return func(ctx *Context) {
		ctx.Write(cache)
	}
}

func (r *Router) getStringServer(str string) func(*Context) {
	b := []byte(str)
	return func(ctx *Context) {
		ctx.Write(b)
	}
}

func (r *Router) getHTMLServer(str HTML) func(*Context) {
	b := []byte(str)
	return func(ctx *Context) {
		if _, err := ctx.HTML().Write(b); err != nil {
			// connection broken
			ctx.Error("", 500)
		}
	}
}

func (r *Router) getJSONServer(str JSON) func(*Context) {
	b := []byte(str)
	return func(ctx *Context) {
		ctx.SetContentType(jsonCTshort)
		ctx.Write(b)
	}
}

func (r *Router) getBytesServer(b []byte) func(*Context) {
	return func(ctx *Context) {
		ctx.Write(b)
	}
}

func (r *Router) getFmtDHandler(v interface{}) func(*Context) {
	const fmtD = "%d"
	return func(ctx *Context) {
		fmt.Fprintf(ctx, fmtD, v)
	}
}

func (r *Router) getFmtFHandler(v interface{}) func(*Context) {
	const fmtF = "%f"
	return func(ctx *Context) {
		fmt.Fprintf(ctx, fmtF, v)
	}
}

// PanicHandler set a handler for unhandled panics
func (r *Router) PanicHandler(panicHandler func(*Context, interface{})) {
	r.initRouter()
	r.router.PanicHandler = panicHandler
}

// NotFound set a handler which is called when no matching route is found
func (r *Router) NotFound(notFoundHandler func(*Context)) {
	r.initRouter()
	r.router.NotFound = notFoundHandler
}

// HandleMethodNotAllowed changes HandleMethodNotAllowed mode in the router
func (r *Router) HandleMethodNotAllowed(newValue bool) (oldValue bool) {
	r.initRouter()
	oldValue = r.router.HandleMethodNotAllowed
	r.router.HandleMethodNotAllowed = newValue
	return
}

// HandleOPTIONS changes HandleOPTIONS mode in the router
func (r *Router) HandleOPTIONS(newValue bool) (oldValue bool) {
	r.initRouter()
	oldValue = r.router.HandleOPTIONS
	r.router.HandleOPTIONS = newValue
	return
}

// HTTP router returns a router instance that work only on HTTP requests
func (r *Router) HTTP() *Router {
	if r.root != nil {
		return r.root.HTTP()
	}
	r.mu.Lock()
	if r.httprouter == nil {
		r.httprouter = &Router{
			router: newRouter(),
			app:    r.app,
			root:   r,
		}
	}
	r.mu.Unlock()

	return r.httprouter
}

// HTTPS router returns a router instance that work only on HTTPS requests
func (r *Router) HTTPS() *Router {
	if r.root != nil {
		return r.root.HTTPS()
	}
	r.mu.Lock()
	if r.httpsrouter == nil {
		r.httpsrouter = &Router{
			router: newRouter(),
			app:    r.app,
			root:   r,
		}
	}
	r.mu.Unlock()

	return r.httpsrouter
}

// ServeFiles serves files from the given file system root.
// The path must end with "/*filepath", files are then served from the local
// path /defined/root/dir/*filepath.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
//     router.ServeFiles("/src/*filepath", "/var/www")
func (r *Router) ServeFiles(path string, rootPath string) {
	r.router.ServeFiles(path, rootPath, nil)
}

// Lookup allows the manual lookup of a method + path combo.
// This is e.g. useful to build a framework around this router.
// If the path was found, it returns the handle function and the path parameter
// values. Otherwise the third return value indicates whether a redirection to
// the same path with an extra / without the trailing slash should be performed.
func (r *Router) Lookup(method, path string, ctx *Context) (RequestHandler, bool) {
	return r.router.Lookup(method, path, ctx)
}

// MethodNotAllowed sets MethodNotAllowed handler
func (r *Router) MethodNotAllowed(c func(ctx *Context)) {
	r.router.MethodNotAllowed = c
}

// Allowed returns Allow header's value used in OPTIONS responses
func (r *Router) Allowed(path, reqMethod string) (allow string) {
	return r.router.Allowed(path, reqMethod)
}

// Handler makes the router implement the fasthttp.ListenAndServe interface.
func (r *Router) Handler() func(*Context) {
	return r.handler
}

func (r *Router) handler(ctx *Context) {
	path := string(ctx.Path())
	method := string(ctx.Method())

	switch ctx.IsTLS() {
	case true:
		if r.httpsrouter != nil {
			handler, rts := r.httpsrouter.router.Lookup(method, path, ctx)
			if handler != nil && r.httpsrouter.handle(path, method, ctx, handler, rts, false) {
				return
			}
			if r.httpsrouter.router.RedirectFixedPath {
				if root, ok := r.httpsrouter.router.Trees[method]; ok && root != nil {
					fixedPath, found := root.FindCaseInsensitivePath(
						CleanPath(path),
						r.httpsrouter.router.RedirectTrailingSlash,
					)

					if found {
						code := redirectCode
						if method != GET {
							code = temporaryRedirectCode
						}
						ctx.SetStatusCode(code)
						ctx.Response.Header.AddBytesV("Location", fixedPath)
						return
					}
				}
			}

			if r.httpsrouter.router.NotFound != nil {
				r.httpsrouter.router.NotFound(ctx)
				return
			}
		}
	case false:
		if r.httprouter != nil {
			handler, rts := r.httprouter.router.Lookup(method, path, ctx)
			if handler != nil && r.httprouter.handle(path, method, ctx, handler, rts, false) {
				return
			}

			if r.httprouter.router.RedirectFixedPath {
				if root, ok := r.httprouter.router.Trees[method]; ok && root != nil {
					fixedPath, found := root.FindCaseInsensitivePath(
						CleanPath(path),
						r.httprouter.router.RedirectTrailingSlash,
					)

					if found {
						code := redirectCode
						if method != GET {
							code = temporaryRedirectCode
						}
						ctx.SetStatusCode(code)
						ctx.Response.Header.AddBytesV("Location", fixedPath)
						return
					}
				}
			}
			if r.httprouter.router.NotFound != nil {
				r.httprouter.router.NotFound(ctx)
				return
			}
		}
	}
	handler, rts := r.router.Lookup(method, path, ctx)
	if r.handle(path, method, ctx, handler, rts, true) {
		return
	}
	if r.router.RedirectFixedPath {
		if root, ok := r.router.Trees[method]; ok && root != nil {
			fixedPath, found := root.FindCaseInsensitivePath(
				CleanPath(path),
				r.router.RedirectTrailingSlash,
			)

			if found && len(fixedPath) > 0 {
				code := redirectCode
				if method != GET {
					code = temporaryRedirectCode
				}
				ctx.SetStatusCode(code)
				ctx.Response.Header.AddBytesV("Location", fixedPath)
				return
			}
		}
	}
	if r.router.NotFound != nil {
		r.router.NotFound(ctx)
		return
	}
	ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
}

func (r *Router) handle(path, method string, ctx *Context, handler func(ctx *Context), redirectTrailingSlashs bool, isRootRouter bool) (handlerFound bool) {
	if r.router.PanicHandler != nil {
		defer r.router.Recv(ctx, nil)
	}
	if f, ok := r.router.StaticHandlers[method][path]; ok {
		f(ctx)
		return true
	}
	if root := r.router.Trees[method]; root != nil {
		if f, tsr := root.GetValue(path, ctx, string(ctx.Method())); f != nil {
			f(ctx)
			return true
		} else if method != CONNECT && path != PathSlash {
			code := redirectCode // Permanent redirect, request with GET method
			if method != GET {
				// Temporary redirect, request with same method
				// As of Go 1.3, Go does not support status code 308.
				code = temporaryRedirectCode
			}

			if tsr && r.router.RedirectTrailingSlash {
				var uri string
				if len(path) > one && path[len(path)-one] == SlashByte {
					uri = path[:len(path)-one]
				} else {
					uri = path + PathSlash
				}
				if uri != emptyString {
					ctx.SetStatusCode(code)
					ctx.Response.Header.Add("Location", uri)
				}
				return false
			}

			// Try to fix the request path
			if r.router.RedirectFixedPath {
				fixedPath, found := root.FindCaseInsensitivePath(
					CleanPath(path),
					r.router.RedirectTrailingSlash,
				)

				if found && len(fixedPath) > 0 {
					queryBuf := ctx.URI().QueryString()
					if len(queryBuf) > zero {
						fixedPath = append(fixedPath, QuestionMark...)
						fixedPath = append(fixedPath, queryBuf...)
					}
					uri := string(fixedPath)
					ctx.SetStatusCode(code)
					ctx.Response.Header.Add("Location", uri)
					return true
				}
			}
		}
	}

	if !isRootRouter {
		return false
	}

	if method == OPTIONS {
		// Handle OPTIONS requests
		if r.router.HandleOPTIONS {
			if allow := r.router.Allowed(path, method); len(allow) > zero {
				ctx.Response.Header.Set(HeaderAllow, allow)
				return true
			}
		}
	} else {
		// Handle 405
		if r.router.HandleMethodNotAllowed {
			if allow := r.router.Allowed(path, method); len(allow) > zero {
				ctx.Response.Header.Set(HeaderAllow, allow)
				if r.router.MethodNotAllowed != nil {
					r.router.MethodNotAllowed(ctx)
				} else {
					ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
					ctx.SetContentTypeBytes(DefaultContentType)
					ctx.SetBodyString(fasthttp.StatusMessage(fasthttp.StatusMethodNotAllowed))
				}
				return true
			}
		}
	}

	return false
}

// Redir sends 301 redirect to the given url
//
// it's equivalent to
//
//     ctx.Redirect(url, 301)
func (r *Router) Redir(route, url string) {
	r.GET(route, func(ctx *Context) {
		ctx.Redirect(route, redirectCode)
	})
}
