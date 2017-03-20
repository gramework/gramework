package gramework

import (
	"os"
	"sync"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

// New App
func New() *App {
	logger := &log.Logger{
		Level:   log.InfoLevel,
		Handler: cli.New(os.Stdout),
	}
	flags := &Flags{
		values: make(map[string]Flag, 0),
	}
	app := &App{
		Flags:                     flags,
		flagsQueue:                flagsToRegister,
		Logger:                    logger,
		domainListLock:            &sync.RWMutex{},
		domains:                   make(map[string]*Router, 0),
		middlewaresMu:             &sync.RWMutex{},
		middlewaresAfterRequestMu: &sync.RWMutex{},
		preMiddlewaresMu:          &sync.RWMutex{},
		middlewares:               make([]func(*Context), 0),
		middlewaresAfterRequest:   make([]func(*Context), 0),
		preMiddlewares:            make([]func(*Context), 0),
	}

	app.defaultRouter = &Router{
		router: newRouter(),
		app:    app,
	}

	return app
}
