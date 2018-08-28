// Copyright 2017-present Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package gramework

import (
	"github.com/apex/log"
	"github.com/google/uuid"
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

		xReqID := ctx.Request.Header.Peek(xRequestID)
		if len(xReqID) > 0 {
			ctx.requestID = string(xReqID)
		} else {
			ctx.requestID = uuid.New().String()
		}

		tracer := ctx.Logger.
			WithFields(log.Fields{
				"package":  "gramework",
				"method":   BytesToString(ctx.Method()),
				"path":     BytesToString(ctx.Path()),
				xRequestID: ctx.requestID,
			})

		if app.defaultRouter.router.PanicHandler != nil || !app.NoDefaultPanicHandler {
			// unfortunately, we can't get rid of that defer
			defer app.defaultRouter.router.Recv(ctx, tracer)
		}

		ctx.Response.Header.Add(xRequestID, ctx.requestID)
		ctx.Logger = ctx.Logger.WithField(xRequestID, ctx.requestID)

		ctx.loadCookies()
		app.preMiddlewaresMu.RLock()
		for k := range app.preMiddlewares {
			app.preMiddlewares[k](ctx)
		}

		app.preMiddlewaresMu.RUnlock()
		if ctx.middlewareKilledReq {
			ctx.saveCookies()
			tracer.
				WithField("status", ctx.Response.StatusCode()).
				Debug("middleware stopped processing")
			return
		}
		ctx.middlewaresShouldStopProcessing = false
		app.middlewaresMu.RLock()
		for k := range app.middlewares {
			app.middlewares[k](ctx)
			if ctx.middlewaresShouldStopProcessing {
				break
			}
		}

		app.middlewaresMu.RUnlock()
		if ctx.middlewareKilledReq {
			ctx.saveCookies()
			tracer.
				WithField("status", ctx.Response.StatusCode()).
				Debug("middleware stopped processing")
			return
		}
		if len(app.domains) > 0 {
			app.domainListLock.RLock()
			if app.domains[string(ctx.URI().Host())] != nil {
				app.domainListLock.RUnlock()
				app.domains[string(ctx.URI().Host())].handler(ctx)
				app.runMiddlewaresAfterRequest(ctx)
				ctx.saveCookies()
				tracer.
					WithField("status", ctx.Response.StatusCode()).
					Debug("request processed")
				return
			}

			app.domainListLock.RUnlock()
			if !app.HandleUnknownDomains {
				ctx.NotFound()
				app.runMiddlewaresAfterRequest(ctx)
				ctx.saveCookies()
				tracer.
					WithField("status", ctx.Response.StatusCode()).
					Debug("request processed")
				return
			}
		}

		app.defaultRouter.handler(ctx)
		app.runMiddlewaresAfterRequest(ctx)
		ctx.saveCookies()
		tracer.
			WithField("status", ctx.Response.StatusCode()).
			Debug("request processed")
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
