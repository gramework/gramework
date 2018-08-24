package gramework

import "github.com/valyala/fasthttp"

func (app *App) copyServer() *fasthttp.Server {
	return &fasthttp.Server{
		Handler:                       app.serverBase.Handler,
		Name:                          app.serverBase.Name,
		Concurrency:                   app.serverBase.Concurrency,
		DisableKeepalive:              app.serverBase.DisableKeepalive,
		ReadBufferSize:                app.serverBase.ReadBufferSize,
		WriteBufferSize:               app.serverBase.WriteBufferSize,
		ReadTimeout:                   app.serverBase.ReadTimeout,
		WriteTimeout:                  app.serverBase.WriteTimeout,
		MaxConnsPerIP:                 app.serverBase.MaxConnsPerIP,
		MaxRequestsPerConn:            app.serverBase.MaxRequestsPerConn,
		MaxKeepaliveDuration:          app.serverBase.MaxKeepaliveDuration,
		MaxRequestBodySize:            app.serverBase.MaxRequestBodySize,
		ReduceMemoryUsage:             app.serverBase.ReduceMemoryUsage,
		GetOnly:                       app.serverBase.GetOnly,
		LogAllErrors:                  app.serverBase.LogAllErrors,
		DisableHeaderNamesNormalizing: app.serverBase.DisableHeaderNamesNormalizing,
		Logger: app.serverBase.Logger,
	}
}
