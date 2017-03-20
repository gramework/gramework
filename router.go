package gramework

import (
	"fmt"

	"time"

	"github.com/kirillDanshin/fasthttprouter"
	"github.com/valyala/fasthttp"
)

// GET registers a handler for a GET request to the given route
func (r *Router) GET(route string, handler interface{}) *Router {
	r.Handle(MethodGET, route, handler)
	return r
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

// Handle registers a new request handle with the given path and method.
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut functions can be used.
// This function is intended for bulk loading and to allow the usage of less frequently used,
// non-standardized or custom methods (e.g. for internal communication with a proxy).
func (r *Router) Handle(method, route string, handler interface{}) *Router {
	r.initRouter()

	switch h := handler.(type) {
	case func(*fasthttp.RequestCtx):
		r.router.Handle(method, route, r.getGrameHandler(h))
	case func(*Context):
		r.router.Handle(method, route, h)
	case func(*fasthttp.RequestCtx) error:
		r.router.Handle(method, route, r.getGrameErrorHandler(h))
	case func(*Context) error:
		r.router.Handle(method, route, r.getErrorHandler(h))
	case string:
		r.router.Handle(method, route, r.getStringServer(h))
	case []byte:
		r.router.Handle(method, route, r.getBytesServer(h))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		r.router.Handle(method, route, r.getFmtDHandler(h))
	case float32, float64:
		r.router.Handle(method, route, r.getFmtFHandler(h))
	default:
		r.app.Logger.Warnf("Unknown handler type: %T, serving fmt.Sprintf(%%v)", h)
		r.router.Handle(method, route, r.getFmtVHandler(h))
	}
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

// NotFound set a handler wich is called when no matching route is found
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

// ServeDir from a given path
func (app *App) ServeDir(path string) func(*Context) {
	return app.ServeDirCustom(path, 0, true, false, nil)
}

// ServeDirCustom gives you ability to serve a dir with custom settings
func (app *App) ServeDirCustom(path string, stripSlashes int, compress bool, generateIndexPages bool, indexNames []string) func(*Context) {
	if indexNames == nil {
		indexNames = []string{}
	}
	fs := &fasthttp.FS{
		Root:                 path,
		IndexNames:           indexNames,
		GenerateIndexPages:   generateIndexPages,
		Compress:             compress,
		CacheDuration:        5 * time.Minute,
		CompressedFileSuffix: ".gz",
	}

	if stripSlashes > 0 {
		fs.PathRewrite = fasthttp.NewPathSlashesStripper(stripSlashes)
	}

	h := fs.NewRequestHandler()
	return func(ctx *Context) {
		h(ctx.RequestCtx)
	}
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
	r.router.ServeFiles(path, rootPath)
}

// Lookup allows the manual lookup of a method + path combo.
// This is e.g. useful to build a framework around this router.
// If the path was found, it returns the handle function and the path parameter
// values. Otherwise the third return value indicates whether a redirection to
// the same path with an extra / without the trailing slash should be performed.
func (r *Router) Lookup(method, path string, ctx *Context) (RequestHandler, bool) {
	return r.router.Lookup(method, path, ctx)
}

// Allowed returns Allow header's value used in OPTIONS responses
func (r *Router) Allowed(path, reqMethod string) (allow string) {
	return r.router.Allowed(path, reqMethod)
}

// Handler makes the router implement the fasthttp.ListenAndServe interface.
func (r *Router) Handler() func(*Context) {
	return func(ctx *Context) {
		path := string(ctx.Path())
		method := string(ctx.Method())
		switch ctx.IsTLS() {
		case true:
			if r.httpsrouter != nil {
				handler, rts := r.httpsrouter.router.Lookup(method, path, ctx)
				if handler != nil && r.httpsrouter.handle(path, method, ctx, handler, rts, false) {
					return
				}
			}
		case false:
			if r.httprouter != nil {
				handler, rts := r.httprouter.router.Lookup(method, path, ctx)
				if handler != nil && r.httprouter.handle(path, method, ctx, handler, rts, false) {
					return
				}
			}
		}
		handler, rts := r.router.Lookup(method, path, ctx)
		if r.handle(path, method, ctx, handler, rts, true) {
			return
		}
		if r.router.NotFound != nil {
			r.router.NotFound(ctx)
			return
		}
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
	}
}

func (r *Router) handle(path, method string, ctx *Context, handler func(ctx *Context), redirectTrailingSlashs bool, isRootRouter bool) (handlerFound bool) {
	if root := r.router.Trees[method]; root != nil {
		if f, tsr := root.GetValue(path, ctx); f != nil {
			f(ctx)
			return true
		} else if method != fasthttprouter.CONNECT && path != fasthttprouter.PathSlash {
			code := redirectCode // Permanent redirect, request with GET method
			if method != fasthttprouter.GET {
				// Temporary redirect, request with same method
				// As of Go 1.3, Go does not support status code 308.
				code = temporaryRedirectCode
			}

			if tsr && r.router.RedirectTrailingSlash {
				var uri string
				if len(path) > one && path[len(path)-one] == fasthttprouter.SlashByte {
					uri = path[:len(path)-one]
				} else {
					uri = path + fasthttprouter.PathSlash
				}
				ctx.Redirect(uri, code)
				return false
			}

			// Try to fix the request path
			if r.router.RedirectFixedPath {
				fixedPath, found := root.FindCaseInsensitivePath(
					fasthttprouter.CleanPath(path),
					r.router.RedirectTrailingSlash,
				)

				if found {
					queryBuf := ctx.URI().QueryString()
					if len(queryBuf) > zero {
						fixedPath = append(fixedPath, fasthttprouter.QuestionMark...)
						fixedPath = append(fixedPath, queryBuf...)
					}
					uri := string(fixedPath)
					ctx.Redirect(uri, code)
					return true
				}
			}
		}
	}

	if !isRootRouter {
		return false
	}

	if method == fasthttprouter.OPTIONS {
		// Handle OPTIONS requests
		if r.router.HandleOPTIONS {
			if allow := r.router.Allowed(path, method); len(allow) > zero {
				ctx.Response.Header.Set(fasthttprouter.HeaderAllow, allow)
				return true
			}
		}
	} else {
		// Handle 405
		if r.router.HandleMethodNotAllowed {
			if allow := r.router.Allowed(path, method); len(allow) > zero {
				ctx.Response.Header.Set(fasthttprouter.HeaderAllow, allow)
				if r.router.MethodNotAllowed != nil {
					r.router.MethodNotAllowed(ctx)
				} else {
					ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
					ctx.SetContentTypeBytes(fasthttprouter.DefaultContentType)
					ctx.SetBodyString(fasthttp.StatusMessage(fasthttp.StatusMethodNotAllowed))
				}
				return true
			}
		}
	}

	return false
}
