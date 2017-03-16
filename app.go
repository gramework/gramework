package gramework

// ToTLSHandler returns handler that redirects user to HTTP scheme
func (app *App) ToTLSHandler() func(*Context) {
	return func(ctx *Context) {
		ctx.ToTLS()
	}
}

func (app *App) HTTP() *Router {
	return app.defaultRouter.HTTP()
}
func (app *App) HTTPS() *Router {
	return app.defaultRouter.HTTPS()
}
