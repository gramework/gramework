// Copyright 2017 Kirill Danshin and Gramework contributors
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

	"github.com/apex/log"
	"github.com/gramework/utils/nocopy"
	"github.com/valyala/fasthttp"
)

type (
	// App represents a gramework app
	App struct {
		defaultRouter             *Router
		domains                   map[string]*Router
		errorHandler              func(func(*fasthttp.RequestCtx) error)
		firewall                  *firewall
		firewallInit              *sync.Once
		Flags                     *Flags
		flagsQueue                []Flag
		Logger                    log.Interface
		name                      string
		Settings                  Settings
		TLSEmails                 []string
		middlewares               []func(*Context)
		middlewaresAfterRequest   []func(*Context)
		preMiddlewares            []func(*Context)
		domainListLock            *sync.RWMutex
		middlewaresAfterRequestMu *sync.RWMutex
		middlewaresMu             *sync.RWMutex
		preMiddlewaresMu          *sync.RWMutex
		EnableFirewall            bool
		flagsRegistered           bool
		HandleUnknownDomains      bool
	}

	// Context is a gramework request context
	Context struct {
		*fasthttp.RequestCtx
		nocopy    nocopy.NoCopy
		Logger    log.Interface
		App       *App
		auth      *Auth
		Cookies   Cookies
		requestID string

		middlewaresShouldStopProcessing bool
	}

	// Cookies handles a typical cookie storage
	Cookies struct {
		Storage map[string]string
		Mu      sync.RWMutex
	}

	// Settings for an App instance
	Settings struct {
		Firewall FirewallSettings
	}

	// FirewallSettings represents a new firewall settings.
	// Internal firewall representation copies this settings
	// atomically.
	FirewallSettings struct {
		// MaxReqPerMin is a max request per minute count
		MaxReqPerMin int64
		// BlockTimeout in seconds
		BlockTimeout int64
	}

	firewall struct {
		// Store a copy of current settings
		MaxReqPerMin        *int64
		BlockTimeout        *int64
		blockListMutex      sync.Mutex
		requestCounterMutex sync.Mutex
		blockList           map[string]int64
		requestCounter      map[string]int64
	}

	// Flags is a flags storage
	Flags struct {
		values map[string]Flag
	}

	// Flag is a flag representation
	Flag struct {
		Name        string
		Description string
		Default     string
		Value       *string
	}

	// Router handles internal handler conversion etc.
	Router struct {
		router      *router
		httprouter  *Router
		httpsrouter *Router
		root        *Router
		app         *App
		mu          sync.RWMutex
		submu       sync.Mutex
	}

	// SubRouter handles subs registration
	// like app.Sub("v1").GET("someRoute", "hi")
	SubRouter struct {
		parent routerable
		prefix string
	}

	routerable interface {
		handleReg(method, route string, handler interface{})
		determineHandler(handler interface{}) func(*Context)
	}

	// RequestHandler describes a standard request handler type
	RequestHandler func(*Context)

	// RequestHandlerErr describes a standard request handler with error returned type
	RequestHandlerErr func(*Context) error

	// Auth is a struct that handles
	// context's basic auth features
	Auth struct {
		login string
		pass  string

		parsed bool
		// if error occurred during parsing,
		// it will be always returned for current
		// context
		err error

		ctx *Context
	}
)
