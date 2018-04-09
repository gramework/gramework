package gramework

import (
	"github.com/valyala/fasthttp"
)

func (app *App) handler() func(*fasthttp.RequestCtx) {
	return func(fhctx *fasthttp.RequestCtx) {
		if app.EnableFirewall {
			app.firewallInit.Do(func() {
				app.initFirewall()
			})
		}
		ctx := app.defaultRouter.initGrameCtx(fhctx)
		if app.EnableFirewall {
			if shouldBeBlocked, _ := app.firewall.NewRequest(ctx); shouldBeBlocked {
				ctx.SetConnectionClose()
				return
			}
		}
		if app.defaultRouter.router.PanicHandler != nil {
			defer app.defaultRouter.router.Recv(ctx)
		}
		ctx.loadCookies()
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
				ctx.saveCookies()
				return
			}
			app.domainListLock.RUnlock()
			if !app.HandleUnknownDomains {
				ctx.NotFound()
				app.runMiddlewaresAfterRequest(ctx)
				ctx.saveCookies()
				return
			}
		}
		app.defaultRouter.Handler()(ctx)
		app.runMiddlewaresAfterRequest(ctx)
		ctx.saveCookies()
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
