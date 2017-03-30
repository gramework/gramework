// Copyright 2013 Julien Schmidt. All rights reserved.
// Copyright (c) 2015-2016, 招牌疯子
// Copyright (c) 2017, Kirill Danshin
// Use of this source code is governed by a BSD-style license that can be found
// in the 3rd-Party License/fasthttprouter file.

package gramework

import (
	"strings"

	"github.com/valyala/fasthttp"
)

var (
	// DefaultContentType cached to minimize memory allocations
	DefaultContentType = []byte("text/plain; charset=utf-8")
	// QuestionMark cached to minimize memory allocations
	QuestionMark = []byte("?")

	// SlashByte cached to minimize memory allocations
	SlashByte = byte('/')
)

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
type router struct {
	Trees map[string]*node

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 307 for all other request methods.
	RedirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 307 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// If enabled, the router automatically replies to OPTIONS requests.
	// Custom OPTIONS handlers take priority over automatic replies.
	HandleOPTIONS bool

	// Configurable http.Handler which is called when no matching route is
	// found. If it is not set, http.NotFound is used.
	NotFound RequestHandler

	// Configurable http.Handler which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
	// The "Allow" header with allowed request methods is set before the handler
	// is called.
	MethodNotAllowed RequestHandler

	// Function to handle panics recovered from http handlers.
	// It should be used to generate a error page and return the http error code
	// 500 (Internal Server Error).
	// The handler can be used to keep your server from crashing because of
	// unrecovered panics.
	PanicHandler func(*Context, interface{})

	cache *cache
}

// newRouter returns a new initialized Router.
// Path auto-correction, including trailing slashes, is enabled by default.
func newRouter() *router {
	r := &router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
		cache: &cache{
			v: make(map[string]*msc, 0),
		},
	}
	go r.cache.maintain()
	return r
}

const (
	// GET method
	GET = "GET"
	// HEAD method
	HEAD = "HEAD"
	// OPTIONS method
	OPTIONS = "OPTIONS"
	// POST method
	POST = "POST"
	// PUT method
	PUT = "PUT"
	// PATCH method
	PATCH = "PATCH"
	// DELETE method
	DELETE = "DELETE"
	// CONNECT method
	CONNECT = "CONNECT"

	// PathAny used to minimize memory allocations
	PathAny = "*"
	// PathSlashAny used to minimize memory allocations
	PathSlashAny = "/*"
	// PathSlash used to minimize memory allocations
	PathSlash = "/"

	// HeaderAllow used to minimize memory allocations
	HeaderAllow = "Allow"
)

// GET is a shortcut for router.Handle("GET", path, handle)
func (r *router) GET(path string, handle RequestHandler) {
	r.Handle(GET, path, handle)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (r *router) HEAD(path string, handle RequestHandler) {
	r.Handle(HEAD, path, handle)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (r *router) OPTIONS(path string, handle RequestHandler) {
	r.Handle(OPTIONS, path, handle)
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (r *router) POST(path string, handle RequestHandler) {
	r.Handle(POST, path, handle)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (r *router) PUT(path string, handle RequestHandler) {
	r.Handle(PUT, path, handle)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (r *router) PATCH(path string, handle RequestHandler) {
	r.Handle(PATCH, path, handle)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (r *router) DELETE(path string, handle RequestHandler) {
	r.Handle(DELETE, path, handle)
}

// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *router) Handle(method, path string, handle RequestHandler) {
	if path[0] != SlashByte {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if r.Trees == nil {
		r.Trees = make(map[string]*node)
	}

	root := r.Trees[method]
	if root == nil {
		root = new(node)
		root.router = r
		r.Trees[method] = root
	}

	root.addRoute(path, handle, r)
}

// ServeFiles serves files from the given file system root.
// The path must end with "/*filepath", files are then served from the local
// path /defined/root/dir/*filepath.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
//     router.ServeFiles("/src/*filepath", "/var/www")
func (r *router) ServeFiles(path string, rootPath string) {
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path must end with /*filepath in path '" + path + "'")
	}
	prefix := path[:len(path)-10]

	fileHandler := fasthttp.FSHandler(rootPath, strings.Count(prefix, PathSlash))

	r.GET(path, func(ctx *Context) {
		fileHandler(ctx.RequestCtx)
	})
}

// Recv used to recover after panic. Called if PanicHandler was set
func (r *router) Recv(ctx *Context) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(ctx, rcv)
	}
}

// Lookup allows the manual lookup of a method + path combo.
// This is e.g. useful to build a framework around this router.
// If the path was found, it returns the handle function and the path parameter
// values. Otherwise the third return value indicates whether a redirection to
// the same path with an extra / without the trailing slash should be performed.
func (r *router) Lookup(method, path string, ctx *Context) (RequestHandler, bool) {
	if root := r.Trees[method]; root != nil {
		return root.GetValue(path, ctx, method)
	}
	return nil, false
}

// Allowed returns Allow header's value used in OPTIONS responses
func (r *router) Allowed(path, reqMethod string) (allow string) {
	if path == PathAny || path == PathSlashAny { // server-wide
		for method := range r.Trees {
			if method == OPTIONS {
				continue
			}

			// add request method to list of allowed methods
			if len(allow) == 0 {
				allow = method
			} else {
				allow += ", " + method
			}
		}
	} else { // specific path
		for method := range r.Trees {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == OPTIONS {
				continue
			}

			handle, _ := r.Trees[method].GetValue(path, nil, reqMethod)
			if handle != nil {
				// add request method to list of allowed methods
				if len(allow) == 0 {
					allow = method
				} else {
					allow += ", " + method
				}
			}
		}
	}
	if len(allow) > 0 {
		allow += ", OPTIONS"
	}
	return
}

// // Handler makes the router implement the fasthttp.ListenAndServe interface.
// func (r *router) Handler(ctx *Context) {
// 	if r.PanicHandler != nil {
// 		defer r.Recv(ctx)
// 	}

// 	path := string(ctx.Path())
// 	method := string(ctx.Method())
// 	if root := r.Trees[method]; root != nil {
// 		if f, tsr := root.GetValue(path, ctx); f != nil {
// 			f(ctx)
// 			return
// 		} else if method != CONNECT && path != PathSlash {
// 			code := 301 // Permanent redirect, request with GET method
// 			if method != GET {
// 				// Temporary redirect, request with same method
// 				// As of Go 1.3, Go does not support status code 308.
// 				code = 307
// 			}

// 			if tsr && r.RedirectTrailingSlash {
// 				var uri string
// 				if len(path) > 1 && path[len(path)-1] == SlashByte {
// 					uri = path[:len(path)-1]
// 				} else {
// 					uri = path + PathSlash
// 				}
// 				ctx.Redirect(uri, code)
// 				return
// 			}

// 			// Try to fix the request path
// 			if r.RedirectFixedPath {
// 				fixedPath, found := root.FindCaseInsensitivePath(
// 					CleanPath(path),
// 					r.RedirectTrailingSlash,
// 				)

// 				if found {
// 					queryBuf := ctx.URI().QueryString()
// 					if len(queryBuf) > 0 {
// 						fixedPath = append(fixedPath, QuestionMark...)
// 						fixedPath = append(fixedPath, queryBuf...)
// 					}
// 					uri := string(fixedPath)
// 					ctx.Redirect(uri, code)
// 					return
// 				}
// 			}
// 		}
// 	}

// 	if method == OPTIONS {
// 		// Handle OPTIONS requests
// 		if r.HandleOPTIONS {
// 			if allow := r.Allowed(path, method); len(allow) > 0 {
// 				ctx.Response.Header.Set(HeaderAllow, allow)
// 				return
// 			}
// 		}
// 	} else {
// 		// Handle 405
// 		if r.HandleMethodNotAllowed {
// 			if allow := r.Allowed(path, method); len(allow) > 0 {
// 				ctx.Response.Header.Set(HeaderAllow, allow)
// 				if r.MethodNotAllowed != nil {
// 					r.MethodNotAllowed(ctx)
// 				} else {
// 					ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
// 					ctx.SetContentTypeBytes(DefaultContentType)
// 					ctx.SetBodyString(fasthttp.StatusMessage(fasthttp.StatusMethodNotAllowed))
// 				}
// 				return
// 			}
// 		}
// 	}

// 	// Handle 404
// 	if r.NotFound != nil {
// 		r.NotFound(ctx)
// 	} else {
// 		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound),
// 			fasthttp.StatusNotFound)
// 	}
// }
