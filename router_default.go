package gramework

// JSON register internal handler that sets json content type
// and serves given handler with GET method
func (app *App) JSON(route string, handler interface{}) *App {
	h := app.defaultRouter.determineHandler(handler)
	app.defaultRouter.Handle(MethodGET, route, jsonHandler(h))

	return app
}

func jsonHandler(h func(*Context)) func(*Context) {
	return func(ctx *Context) {
		ctx.SetContentType(jsonCT)
		h(ctx)
	}
}

// GET registers a handler for a GET request to the given route
func (app *App) GET(route string, handler interface{}) *App {
	app.defaultRouter.Handle(MethodGET, route, handler)
	return app
}

// DELETE registers a handler for a DELETE request to the given route
func (app *App) DELETE(route string, handler interface{}) *App {
	app.defaultRouter.Handle(MethodDELETE, route, handler)
	return app
}

// HEAD registers a handler for a HEAD request to the given route
func (app *App) HEAD(route string, handler interface{}) *App {
	app.defaultRouter.Handle(MethodHEAD, route, handler)
	return app
}

// OPTIONS registers a handler for a OPTIONS request to the given route
func (app *App) OPTIONS(route string, handler interface{}) *App {
	app.defaultRouter.Handle(MethodOPTIONS, route, handler)
	return app
}

// PUT registers a handler for a PUT request to the given route
func (app *App) PUT(route string, handler interface{}) *App {
	app.defaultRouter.Handle(MethodPUT, route, handler)
	return app
}

// POST registers a handler for a POST request to the given route
func (app *App) POST(route string, handler interface{}) *App {
	app.defaultRouter.Handle(MethodPOST, route, handler)
	return app
}

// PATCH registers a handler for a PATCH request to the given route
func (app *App) PATCH(route string, handler interface{}) *App {
	app.defaultRouter.Handle(MethodPATCH, route, handler)
	return app
}

// Handle registers a new request handle with the given path and method.
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut functions can be used.
// This function is intended for bulk loading and to allow the usage of less frequently used,
// non-standardized or custom methods (e.g. for internal communication with a proxy).
func (app *App) Handle(method, route string, handler interface{}) *App {
	app.defaultRouter.Handle(method, route, handler)
	return app
}

// PanicHandler set a handler for unhandled panics
func (app *App) PanicHandler(panicHandler func(*Context, interface{})) *App {
	app.defaultRouter.PanicHandler(panicHandler)
	return app
}

// NotFound set a handler which is called when no matching route is found
func (app *App) NotFound(notFoundHandler func(*Context)) *App {
	app.defaultRouter.NotFound(notFoundHandler)
	return app
}

// ServeFile serves a file on a given route
func (app *App) ServeFile(route, file string) *Router {
	return app.defaultRouter.ServeFile(route, file)
}

// SPAIndex serves an index file on any unregistered route
func (app *App) SPAIndex(path string) *Router {
	return app.defaultRouter.SPAIndex(path)
}

// HandleMethodNotAllowed changes HandleMethodNotAllowed mode in the router
func (app *App) HandleMethodNotAllowed(newValue bool) (oldValue bool) {
	return app.defaultRouter.HandleMethodNotAllowed(newValue)
}

// HandleOPTIONS changes HandleOPTIONS mode in the router
func (app *App) HandleOPTIONS(newValue bool) (oldValue bool) {
	return app.defaultRouter.HandleOPTIONS(newValue)
}

// Sub let you quickly register subroutes with given prefix
// like app.Sub("v1").Sub("users").GET("view/:id", "hi").DELETE("delete/:id", "hi"),
// that give you /v1/users/view/:id and /v1/users/delete/:id registered
func (app *App) Sub(path string) *SubRouter {
	return app.defaultRouter.Sub(path)
}
