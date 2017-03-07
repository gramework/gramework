package gramework

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

// New App
func New() *App {
	return &App{
		Flags: &Flags{
			values: make(map[string]Flag, 0),
		},
		flagsQueue: flagsToRegister,
		Logger: &log.Logger{
			Level:   log.InfoLevel,
			Handler: cli.New(os.Stdout),
		},
	}
}
