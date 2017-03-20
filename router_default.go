package gramework

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

// NotFound set a handler wich is called when no matching route is found
func (app *App) NotFound(notFoundHandler func(*Context)) *App {
	app.defaultRouter.NotFound(notFoundHandler)
	return app
}

// HandleMethodNotAllowed changes HandleMethodNotAllowed mode in the router
func (app *App) HandleMethodNotAllowed(newValue bool) (oldValue bool) {
	return app.defaultRouter.HandleMethodNotAllowed(newValue)
}

// HandleOPTIONS changes HandleOPTIONS mode in the router
func (app *App) HandleOPTIONS(newValue bool) (oldValue bool) {
	return app.defaultRouter.HandleOPTIONS(newValue)
}
