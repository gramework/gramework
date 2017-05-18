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
		defaultRouter  *Router
		errorHandler   func(func(*fasthttp.RequestCtx) error)
		EnableFirewall bool
		firewall       *firewall
		firewallInit   *sync.Once

		Logger    log.Interface
		TLSEmails []string
		Settings  Settings

		HandleUnknownDomains bool
		domains              map[string]*Router

		Flags           *Flags
		flagsRegistered bool
		flagsQueue      []Flag

		domainListLock *sync.RWMutex

		preMiddlewares            []func(*Context)
		middlewares               []func(*Context)
		middlewaresAfterRequest   []func(*Context)
		middlewaresMu             *sync.RWMutex
		preMiddlewaresMu          *sync.RWMutex
		middlewaresAfterRequestMu *sync.RWMutex
	}

	// Context is a gramework request context
	Context struct {
		*fasthttp.RequestCtx
		nocopy  nocopy.NoCopy
		Logger  log.Interface
		App     *App
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
		MaxReqPerMin *int64
		BlockTimeout *int64

		blockList      map[string]int64
		blockListMutex sync.Mutex

		requestCounter      map[string]int64
		requestCounterMutex sync.Mutex
	}

	// Flags is a flags storage
	Flags struct {
		values map[string]Flag
	}

	// Flag is a flag representation
	Flag struct {
		Name        string
		Description string
		Value       *string
		Default     string
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
)
