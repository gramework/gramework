package gramework

import (
	"net"
)

// Behind is an interface that allow you to parse provider-dependent
// headers in a compatible way.
//
// WARNING: this interface is currently in a WIP state and unfrozen.
// If you want to support your provider, please open an issue or PR
// at https://github.com/gramework/gramework.
// In that way we can ensure that you have always working implementation.
type Behind interface {
	RemoteIP(ctx *Context) net.IP
	RemoteAddr(ctx *Context) net.Addr
}

func (ctx *Context) RemoteIP() net.IP {
	if ctx.App.behind != nil {
		return ctx.App.behind.RemoteIP(ctx)
	}

	return ctx.RequestCtx.RemoteIP()
}

func (ctx *Context) RemoteAddr() net.Addr {
	if ctx.App.behind != nil {
		return ctx.App.behind.RemoteAddr(ctx)
	}

	return ctx.RequestCtx.RemoteAddr()
}

type internalBehindActivationHook interface {
	OnAppActivation(*App)
}

func (app *App) Behind(u Behind) {
	if hook, ok := u.(internalBehindActivationHook); ok {
		hook.OnAppActivation(app)
	}

	app.behind = u
}
