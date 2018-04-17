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
		newRouter                 func() RouterIface
	}

	// Context is a gramework request context
	Context struct {
		*fasthttp.RequestCtx
		nocopy  nocopy.NoCopy
		Logger  log.Interface
		App     *App
		auth    *Auth
		Cookies Cookies

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
		router      RouterIface
		httprouter  *Router
		httpsrouter *Router
		root        *Router
		app         *App
		mu          sync.RWMutex
		submu       sync.Mutex
	}

	RouterIface interface {
		// MethodNotAllowed sets MethodNotAllowed handler
		MethodNotAllowed(c func(ctx *Context))

		// PanicHandler set a handler for unhandled panics
		PanicHandler(panicHandler func(*Context, interface{}))

		// Recv used to recover after panic. Called if PanicHandler was set
		Recv(ctx *Context)

		// Handle registers a new request handle with the given path and method.
		//
		// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
		// functions can be used.
		//
		// This function is intended for bulk loading and to allow the usage of less
		// frequently used, non-standardized or custom methods (e.g. for internal
		// communication with a proxy).
		Handle(method, path string, handle RequestHandler)

		// HandleMethodNotAllowed changes HandleMethodNotAllowed mode in the router
		HandleMethodNotAllowed(newValue bool) (oldValue bool)

		// HandleOPTIONS changes HandleOPTIONS mode in the router
		HandleOPTIONS(newValue bool) (oldValue bool)

		// SetNotFound set a handler which is called when no matching route is found
		SetNotFound(notFoundHandler func(*Context))

		// NotFound triggers a handler when no matching route is found
		NotFound(*Context)

		// ServeFiles serves files from the given file system root.
		// The path must end with "/*filepath", files are then served from the local
		// path /defined/root/dir/*filepath.
		// For example if root is "/etc" and *filepath is "passwd", the local file
		// "/etc/passwd" would be served.
		// Internally a http.FileServer is used, therefore http.NotFound is used instead
		// of the Router's NotFound handler.
		//     router.ServeFiles("/src/*filepath", "/var/www")
		ServeFiles(path string, rootPath string)

		// Lookup allows the manual lookup of a method + path combo.
		// This is e.g. useful to build a framework around this router.
		// If the path was found, it returns the handle function and the path parameter
		// values. Otherwise the third return value indicates whether a redirection to
		// the same path with an extra / without the trailing slash should be performed.
		Lookup(method, path string, ctx *Context) (RequestHandler, bool)

		// Allowed returns Allow header's value used in OPTIONS responses
		Allowed(path, reqMethod string) (allow string)

		// Process is a callback that will be called in main router's Handle() method
		// on each request.
		Process(method, path string, ctx *Context) (found bool)
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
