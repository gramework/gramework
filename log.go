package gramework

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/valyala/fasthttp"
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

// FastHTTPLoggerAdapter  Adapter for passing apex/log used as gramework Logger into fasthttp
type FastHTTPLoggerAdapter struct {
	apexLogger log.Interface
	fasthttp.Logger
}

// NewFastHTTPLoggerAdapter create new *FastHTTPLoggerAdapter
func NewFastHTTPLoggerAdapter(logger *log.Interface) (fasthttplogger *FastHTTPLoggerAdapter) {
	fasthttplogger = &FastHTTPLoggerAdapter{
		apexLogger: *logger,
	}
	return fasthttplogger
}

//Printf show message only if set app.Logger.Level = apex/log.DebugLevel
func (l *FastHTTPLoggerAdapter) Printf(msg string, v ...interface{}) {
	l.apexLogger.Debugf(msg, v...)
}
