package gramework

// ToTLSHandler returns handler that redirects user to HTTP scheme
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
