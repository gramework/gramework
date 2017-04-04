package gramework

import "fmt"

// GET registers a handler for a GET request to the given route
func (r *SubRouter) GET(route string, handler interface{}) *SubRouter {
	route = r.prefixedRoute(route)
	if r.parent != nil {
		r.parent.handleReg(MethodGET, route, handler)
	}
	return r
}

// DELETE registers a handler for a DELETE request to the given route
func (r *SubRouter) DELETE(route string, handler interface{}) *SubRouter {
	route = r.prefixedRoute(route)
	if r.parent != nil {
		r.parent.handleReg(MethodDELETE, route, handler)
	}
	return r
}

// HEAD registers a handler for a HEAD request to the given route
func (r *SubRouter) HEAD(route string, handler interface{}) *SubRouter {
	route = r.prefixedRoute(route)
	if r.parent != nil {
		r.parent.handleReg(MethodHEAD, route, handler)
	}
	return r
}

// ServeFile serves a file on a given route
func (r *SubRouter) ServeFile(route, file string) *SubRouter {
	route = r.prefixedRoute(route)
	r.parent.handleReg(MethodGET, route, func(ctx *Context) {
		ctx.SendFile(file)
	})
	return r
}

// OPTIONS registers a handler for a OPTIONS request to the given route
func (r *SubRouter) OPTIONS(route string, handler interface{}) *SubRouter {
	route = r.prefixedRoute(route)
	if r.parent != nil {
		r.parent.handleReg(MethodOPTIONS, route, handler)
	}
	return r
}

// PUT registers a handler for a PUT request to the given route
func (r *SubRouter) PUT(route string, handler interface{}) *SubRouter {
	route = r.prefixedRoute(route)
	if r.parent != nil {
		r.parent.handleReg(MethodPUT, route, handler)
	}
	return r
}

// POST registers a handler for a POST request to the given route
func (r *SubRouter) POST(route string, handler interface{}) *SubRouter {
	route = r.prefixedRoute(route)
	if r.parent != nil {
		r.parent.handleReg(MethodPOST, route, handler)
	}
	return r
}

// PATCH registers a handler for a PATCH request to the given route
func (r *SubRouter) PATCH(route string, handler interface{}) *SubRouter {
	route = r.prefixedRoute(route)
	if r.parent != nil {
		r.parent.handleReg(MethodPATCH, route, handler)
	}
	return r
}

func (r *SubRouter) handleReg(method, route string, handler interface{}) {
	r.parent.handleReg(method, route, handler)
}

// Sub let you quickly register subroutes with given prefix
// like app.Sub("v1").Sub("users").GET("view/:id", "hi").DELETE("delete/:id", "hi"),
// that give you /v1/users/view/:id and /v1/users/delete/:id registered
func (r *SubRouter) Sub(path string) *SubRouter {
	return &SubRouter{
		parent: r,
		prefix: r.prefixedRoute(path),
	}
}

func (r *SubRouter) prefixedRoute(route string) string {
	if r.prefix[len(r.prefix)-1] != '/' && route[0] != '/' {
		return fmt.Sprintf("%s/%s", r.prefix, route)
	}
	return fmt.Sprintf("%s%s", r.prefix, route)
}
