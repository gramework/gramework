package gramework

import (
	"github.com/apex/log"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type (
	// App represents a gramework app
	App struct {
		router       *fasthttprouter.Router
		errorHandler func(func(*fasthttp.RequestCtx) error)
		Logger       log.Interface
	}

	// Context is a gramework request context
	Context struct {
		*fasthttp.RequestCtx
		Logger log.Interface
	}
)
