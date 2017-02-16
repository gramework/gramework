package gramework

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

// Logger handles default logger
var Logger = &log.Logger{
	Level:   log.InfoLevel,
	Handler: cli.New(os.Stdout),
}

// Errorf logs an error using default logger
func Errorf(msg string, v ...interface{}) {
	Logger.Errorf(msg, v...)
}
