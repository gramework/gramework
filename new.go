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

import (
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/microcosm-cc/bluemonday"
	"github.com/valyala/fasthttp"
)

var defaultMaxHackAttempts int32 = 5

// New App
func New(opts ...func(*App)) *App {
	logger := Logger
	internalLog = func() *log.Entry {
		Logger.Level = log.DebugLevel
		if !enableDebug {
			Logger.Level = log.InfoLevel
		}

		return Logger.WithField("package", "gramework")
	}()
	flags := &Flags{
		values: make(map[string]Flag),
	}
	defFWLimit := int64(-1)
	defBlockTimeout := int64(-1)
	maxHackAttempts := defaultMaxHackAttempts
	app := &App{
		Flags:          flags,
		flagsQueue:     flagsToRegister,
		Logger:         logger,
		name:           DefaultAppName,
		domainListLock: new(sync.RWMutex),
		firewall: &firewall{
			blockList:      make(map[string]int64),
			MaxReqPerMin:   &defFWLimit,
			BlockTimeout:   &defBlockTimeout,
			requestCounter: make(map[string]int64),
		},
		firewallInit:              new(sync.Once),
		domains:                   make(map[string]*Router),
		middlewaresMu:             new(sync.RWMutex),
		middlewaresAfterRequestMu: new(sync.RWMutex),
		preMiddlewaresMu:          new(sync.RWMutex),
		middlewares:               make([]func(*Context), 0),
		middlewaresAfterRequest:   make([]func(*Context), 0),
		preMiddlewares:            make([]func(*Context), 0),
		seed:                      uintptr(time.Now().Nanosecond()),
		maxHackAttempts:           &maxHackAttempts,
		runningServersMu:          new(sync.Mutex),
		internalLog:               internalLog,
		cookieExpire:              6 * time.Hour,
		cookiePath:                defaultCookiePath,

		sanitizerPolicy: bluemonday.StrictPolicy(),
	}

	for _, opt := range opts {
		opt(app)
	}

	if app.serverBase == nil {
		app.serverBase = newDefaultServerBaseFor(app)
	}
	// avoid race condition then OptUseServer becomes before OptAppName
	// or OptUseCustomLogger becomes before OptUseServer
	app.serverBase.Name = app.name
	app.serverBase.Logger = NewFastHTTPLoggerAdapter(&app.Logger)

	app.defaultRouter = &Router{
		router: newRouter(),
		app:    app,
	}

	return app
}

// Common code used with `gramework.New()` itself and `Opt*` functions.
func newDefaultServerBaseFor(app *App) *fasthttp.Server {
	return &fasthttp.Server{
		Handler: app.handler(),
		Logger:  NewFastHTTPLoggerAdapter(&app.Logger),
	}
}
