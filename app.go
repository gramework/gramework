package gramework

// ToTLSHandler returns handler that redirects user to HTTP scheme
func (app *App) ToTLSHandler() func(*Context) {
	return func(ctx *Context) {
		ctx.ToTLS()
	}
}
