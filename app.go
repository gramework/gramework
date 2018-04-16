// Copyright 2017 Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package gramework

// SetName for the server
func (app *App) SetName(name string) {
	app.name = name
}

// ToTLSHandler returns handler that redirects user to HTTPS scheme
func (app *App) ToTLSHandler() func(*Context) {
	return func(ctx *Context) {
		ctx.ToTLS()
	}
}

// HTTP returns HTTP-only router
func (app *App) HTTP() *Router {
	return app.defaultRouter.HTTP()
}

// HTTPS returns HTTPS-only router
func (app *App) HTTPS() *Router {
	return app.defaultRouter.HTTPS()
}

// MethodNotAllowed sets MethodNotAllowed handler
func (app *App) MethodNotAllowed(c func(ctx *Context)) *App {
	app.defaultRouter.router.MethodNotAllowed = c
	return app
}

// Redir sends 301 redirect to the given url
//
// it's equivalent to
//
//     ctx.Redirect(url, 301)
func (app *App) Redir(url string) func(*Context) {
	return func(ctx *Context) {
		ctx.Redirect(url, redirectCode)
	}
}

// Forbidden send 403 Forbidden error
func (app *App) Forbidden(ctx *Context) {
	ctx.Forbidden()
}
