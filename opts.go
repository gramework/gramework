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
	"errors"

	"github.com/apex/log"
	"github.com/valyala/fasthttp"
)

// OptAppName sets app.name and app.serverBase.Name
func OptAppName(n string) func(*App) {
	return func(app *App) {
		assertAppNotNill(app)
		app.name = n
	}
}

// OptUseCustomLogger allows use custom preconfigured Apex logger.
// For exmaple with custom Handler.
func OptUseCustomLogger(logger *log.Logger) func(*App) {
	return func(app *App) {
		assertAppNotNill(app)
		app.Logger = logger
	}
}

// OptUseServer sets fasthttp.Server instance to use
func OptUseServer(s *fasthttp.Server) func(*App) {
	return func(app *App) {
		assertAppNotNill(app)
		if s == nil {
			panic(errors.New("cannot set nil as app server instance"))
		}
		app.serverBase = s
		app.serverBase.Handler = app.handler()
	}
}

// OptMaxRequestBodySize sets new MaxRequestBodySize in the server used at the execution time.
// All OptUseServer will overwrite this setting 'case OptUseServer replaces the whole server instance
// with a new one.
func OptMaxRequestBodySize(size int) func(*App) {
	return func(app *App) {
		assertAppNotNill(app)
		if app.serverBase == nil {
			app.serverBase = newDefaultServerBaseFor(app)
		}
		app.serverBase.MaxRequestBodySize = size
	}
}

// OptKeepHijackedConns sets new KeepHijackedConns in the server used at the execution time.
// All OptUseServer will overwrite this setting 'case OptUseServer replaces the whole server instance
// with a new one.
func OptKeepHijackedConns(keep bool) func(*App) {
	return func(app *App) {
		assertAppNotNill(app)
		if app.serverBase == nil {
			app.serverBase = newDefaultServerBaseFor(app)
		}
		app.serverBase.KeepHijackedConns = keep
	}
}

func assertAppNotNill(app *App) {
	if app == nil {
		panic(errors.New("option can be implemented only to already creaded app object not to nil"))
	}
}
