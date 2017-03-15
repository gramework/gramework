package gramework

import (
	"github.com/kirillDanshin/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func (r *Router) getErrorHandler(h func(*fasthttp.RequestCtx) error) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		if err := h(ctx); err != nil {
			ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		}
	}
}

func (r *Router) getGrameHandler(h func(*Context)) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		h(r.initGrameCtx(ctx))
	}
}

func (r *Router) getGrameErrorHandler(h func(*Context) error) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		if err := h(r.initGrameCtx(ctx)); err != nil {
			r.app.Logger.WithField("url", ctx.URI()).Errorf("Error occurred: %s", err)
			ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		}
	}
}

func (r *Router) initGrameCtx(ctx *fasthttp.RequestCtx) *Context {
	return &Context{
		Logger:     r.app.Logger,
		RequestCtx: ctx,
	}
}

func (r *Router) initRouter() {
	if r.router == nil {
		r.router = fasthttprouter.New()
	}
}
