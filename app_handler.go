package gramework

import (
	"github.com/valyala/fasthttp"
)

func (app *App) handler() func(*fasthttp.RequestCtx) {
	return func(fhctx *fasthttp.RequestCtx) {
		ctx := app.defaultRouter.initGrameCtx(fhctx)
		app.preMiddlewaresMu.RLock()
		for k := range app.preMiddlewares {
			app.preMiddlewares[k](ctx)
		}
		app.preMiddlewaresMu.RUnlock()
		ctx.middlewaresShouldStopProcessing = false
		app.middlewaresMu.RLock()
		for k := range app.middlewares {
			app.middlewares[k](ctx)
			if ctx.middlewaresShouldStopProcessing {
				break
			}
		}
		app.middlewaresMu.RUnlock()
		if len(app.domains) > 0 {
			d := string(ctx.URI().Host())
			app.domainListLock.RLock()
			if app.domains[d] != nil {
				app.domainListLock.RUnlock()
				app.domains[d].Handler()(ctx)
				app.runMiddlewaresAfterRequest(ctx)
				return
			}
			app.domainListLock.RUnlock()
			if !app.HandleUnknownDomains {
				ctx.NotFound()
				app.runMiddlewaresAfterRequest(ctx)
				return
			}
		}
		app.defaultRouter.Handler()(ctx)
		app.runMiddlewaresAfterRequest(ctx)
	}
}

func (app *App) runMiddlewaresAfterRequest(ctx *Context) {
	app.middlewaresAfterRequestMu.RLock()
	for k := range app.middlewaresAfterRequest {
		app.middlewaresAfterRequest[k](ctx)
		if ctx.middlewaresShouldStopProcessing {
			break
		}
	}
	app.middlewaresAfterRequestMu.RUnlock()
}
