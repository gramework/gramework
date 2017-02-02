package gramework

import "github.com/valyala/fasthttp"

// GET registers a handler for a GET request to the given route
func (app *App) GET(route string, handler interface{}) {
	app.Handle(MethodGET, route, handler)
}

// DELETE registers a handler for a DELETE request to the given route
func (app *App) DELETE(route string, handler interface{}) {
	app.Handle(MethodDELETE, route, handler)
}

// HEAD registers a handler for a HEAD request to the given route
func (app *App) HEAD(route string, handler interface{}) {
	app.Handle(MethodHEAD, route, handler)
}

// OPTIONS registers a handler for a OPTIONS request to the given route
func (app *App) OPTIONS(route string, handler interface{}) {
	app.Handle(MethodOPTIONS, route, handler)
}

// PUT registers a handler for a PUT request to the given route
func (app *App) PUT(route string, handler interface{}) {
	app.Handle(MethodPUT, route, handler)
}

// POST registers a handler for a POST request to the given route
func (app *App) POST(route string, handler interface{}) {
	app.Handle(MethodPOST, route, handler)
}

// PATCH registers a handler for a PATCH request to the given route
func (app *App) PATCH(route string, handler interface{}) {
	app.Handle(MethodPATCH, route, handler)
}

// Handle registers a new request handle with the given path and method.
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut functions can be used.
// This function is intended for bulk loading and to allow the usage of less frequently used,
// non-standardized or custom methods (e.g. for internal communication with a proxy).
func (app *App) Handle(method, route string, handler interface{}) {
	app.initRouter()

	switch h := handler.(type) {
	case func(*fasthttp.RequestCtx):
		app.router.Handle(method, route, h)
	case func(*Context):
		app.router.Handle(method, route, app.getGrameHandler(h))
	case func(*Context) error:
		app.router.Handle(method, route, app.getGrameErrorHandler(h))
	case func(*fasthttp.RequestCtx) error:
		app.router.Handle(method, route, app.getErrorHandler(h))
	}
}

// PanicHandler set a handler for unhandled panics
func (app *App) PanicHandler(panicHandler func(*fasthttp.RequestCtx, interface{})) {
	app.initRouter()
	app.router.PanicHandler = panicHandler
}

// NotFound set a handler wich is called when no matching route is found
func (app *App) NotFound(notFoundHandler func(*fasthttp.RequestCtx)) {
	app.initRouter()
	app.router.NotFound = notFoundHandler
}

// HandleMethodNotAllowed changes HandleMethodNotAllowed mode in the router
func (app *App) HandleMethodNotAllowed(newValue bool) (oldValue bool) {
	app.initRouter()
	oldValue = app.router.HandleMethodNotAllowed
	app.router.HandleMethodNotAllowed = newValue
	return
}

// HandleOPTIONS changes HandleOPTIONS mode in the router
func (app *App) HandleOPTIONS(newValue bool) (oldValue bool) {
	app.initRouter()
	oldValue = app.router.HandleOPTIONS
	app.router.HandleOPTIONS = newValue
	return
}
