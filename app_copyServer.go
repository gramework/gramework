// Copyright 2017-present Kirill Danshin and Gramework contributors
// Copyright 2019-present Highload LTD (UK CN: 11893420)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

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
		Logger:                        app.serverBase.Logger,
		KeepHijackedConns:             app.serverBase.KeepHijackedConns,
	}
}
