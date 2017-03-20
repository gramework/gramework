package gramework

import "github.com/valyala/fasthttp"

func (r *Router) getErrorHandler(h func(*Context) error) func(*Context) {
	return func(ctx *Context) {
		if err := h(ctx); err != nil {
			r.app.Logger.WithField("url", ctx.URI()).Errorf("Error occurred: %s", err)
			ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		}
	}
}

func (r *Router) getGrameHandler(h func(*fasthttp.RequestCtx)) func(*Context) {
	return func(ctx *Context) {
		if ctx != nil {
			h(ctx.RequestCtx)
			return
		}
		h(&fasthttp.RequestCtx{})
	}
}

func (r *Router) getGrameErrorHandler(h func(*fasthttp.RequestCtx) error) func(*Context) {
	return func(ctx *Context) {
		if err := h(ctx.RequestCtx); err != nil {
			r.app.Logger.WithField("url", ctx.URI()).Errorf("Error occurred: %s", err)
			ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		}
	}
}

func (r *Router) initGrameCtx(ctx *fasthttp.RequestCtx) *Context {
	return &Context{
		Logger:     r.app.Logger,
		RequestCtx: ctx,
		App:        r.app,
	}
}

func (r *Router) initRouter() {
	if r.router == nil {
		r.router = newRouter()
	}
}
