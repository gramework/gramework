package gramework

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func (app *App) getErrorHandler(h func(*fasthttp.RequestCtx) error) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		if err := h(ctx); err != nil {
			ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		}
	}
}

func (app *App) getGrameHandler(h func(*Context)) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		h(app.initGrameCtx(ctx))
	}
}

func (app *App) getGrameErrorHandler(h func(*Context) error) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		if err := h(app.initGrameCtx(ctx)); err != nil {
			app.Logger.WithField("url", ctx.URI()).Errorf("Error occurred: %s", err)
			ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		}
	}
}

func (app *App) initGrameCtx(ctx *fasthttp.RequestCtx) *Context {
	return &Context{
		Logger:     app.Logger,
		RequestCtx: ctx,
	}
}

func (app *App) initRouter() {
	if app.router == nil {
		app.router = fasthttprouter.New()
	}
}
