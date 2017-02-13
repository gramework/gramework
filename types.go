package gramework

import (
	"sync"

	"github.com/apex/log"
	"github.com/buaazp/fasthttprouter"
	"github.com/gramework/utils/nocopy"
	"github.com/valyala/fasthttp"
)

type (
	// App represents a gramework app
	App struct {
		router       *fasthttprouter.Router
		errorHandler func(func(*fasthttp.RequestCtx) error)
		firewall     *firewall

		Logger    log.Interface
		TLSEmails []string
		Settings  Settings
	}

	// Context is a gramework request context
	Context struct {
		*fasthttp.RequestCtx
		Logger log.Interface
		nocopy nocopy.NoCopy
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
)
