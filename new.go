package gramework

import (
	"os"
	"sync"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/kirillDanshin/fasthttprouter"
)

// New App
func New() *App {
	app := &App{
		Flags: &Flags{
			values: make(map[string]Flag, 0),
		},
		flagsQueue: flagsToRegister,
		Logger: &log.Logger{
			Level:   log.InfoLevel,
			Handler: cli.New(os.Stdout),
		},
		domainListLock: &sync.RWMutex{},
		domains:        make(map[string]*Router, 0),
	}

	app.defaultRouter = &Router{
		router: fasthttprouter.New(),
		app:    app,
	}

	return app
}
