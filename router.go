package gramework

import (
	"fmt"

	"time"

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
		r.router.Handle(method, route, h)
	case func(*Context):
		r.router.Handle(method, route, r.getGrameHandler(h))
	case func(*Context) error:
		r.router.Handle(method, route, r.getGrameErrorHandler(h))
	case func(*fasthttp.RequestCtx) error:
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

func (r *Router) getFmtVHandler(v interface{}) func(*fasthttp.RequestCtx) {
	cache := []byte(fmt.Sprintf("%v", v))
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Write(cache)
	}
}

func (r *Router) getStringServer(str string) func(*fasthttp.RequestCtx) {
	b := []byte(str)
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Write(b)
	}
}

func (r *Router) getBytesServer(b []byte) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Write(b)
	}
}

func (r *Router) getFmtDHandler(v interface{}) func(*fasthttp.RequestCtx) {
	const fmtD = "%d"
	return func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, fmtD, v)
	}
}

func (r *Router) getFmtFHandler(v interface{}) func(*fasthttp.RequestCtx) {
	const fmtF = "%f"
	return func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, fmtF, v)
	}
}

// PanicHandler set a handler for unhandled panics
func (r *Router) PanicHandler(panicHandler func(*fasthttp.RequestCtx, interface{})) {
	r.initRouter()
	r.router.PanicHandler = panicHandler
}

// NotFound set a handler wich is called when no matching route is found
func (r *Router) NotFound(notFoundHandler func(*fasthttp.RequestCtx)) {
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
func (r *Router) ServeDir(path string) func(*fasthttp.RequestCtx) {
	return r.ServeDirCustom(path, 0, true, false, nil)
}

// ServeDirCustom gives you ability to serve a dir with custom settings
func (r *Router) ServeDirCustom(path string, stripSlashes int, compress bool, generateIndexPages bool, indexNames []string) func(*fasthttp.RequestCtx) {
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
	return func(ctx *fasthttp.RequestCtx) {
		h(ctx)
	}
}
