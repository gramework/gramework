package gramework

import "github.com/valyala/fasthttp"

func (app *App) copyServer() *fasthttp.Server {
	return &fasthttp.Server{
		Handler:                       app.server.Handler,
		Name:                          app.server.Name,
		Concurrency:                   app.server.Concurrency,
		DisableKeepalive:              app.server.DisableKeepalive,
		ReadBufferSize:                app.server.ReadBufferSize,
		WriteBufferSize:               app.server.WriteBufferSize,
		ReadTimeout:                   app.server.ReadTimeout,
		WriteTimeout:                  app.server.WriteTimeout,
		MaxConnsPerIP:                 app.server.MaxConnsPerIP,
		MaxRequestsPerConn:            app.server.MaxRequestsPerConn,
		MaxKeepaliveDuration:          app.server.MaxKeepaliveDuration,
		MaxRequestBodySize:            app.server.MaxRequestBodySize,
		ReduceMemoryUsage:             app.server.ReduceMemoryUsage,
		GetOnly:                       app.server.GetOnly,
		LogAllErrors:                  app.server.LogAllErrors,
		DisableHeaderNamesNormalizing: app.server.DisableHeaderNamesNormalizing,
		Logger: app.server.Logger,
	}
}
