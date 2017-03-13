package gramework

import (
	"github.com/valyala/fasthttp"
)

func (app *App) handler() func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		if len(app.domains) > 0 {
			d := string(ctx.URI().Host())
			app.domainListLock.RLock()
			if app.domains[d] != nil {
				app.domainListLock.RUnlock()
				app.domains[d].router.Handler(ctx)
				return
			}
			app.domainListLock.RUnlock()
			if !app.HandleUnknownDomains {
				ctx.NotFound()
				return
			}
		}
		app.defaultRouter.router.Handler(ctx)
	}
}
