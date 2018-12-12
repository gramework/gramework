// Copyright 2017-present Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package gramework

import (
	"os"
	"strings"
	"sync/atomic"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/valyala/fasthttp"
)

var enableDebug = true

var currentEnvironment *int32

// Environment defines which environment gramework application runs in.
// It may be useful in various cases.
type Environment int32

const (
	// DEV is the default environment
	DEV Environment = iota
	// STAGE envoronment works just like prod environment,
	// but with detailed logs
	STAGE
	// PROD environment itself
	PROD
)

func (e Environment) String() string {
	switch e {
	case DEV:
		return "DEV"
	case STAGE:
		return "STAGE"
	case PROD:
		return "PROD"
	default:
		return "<unknown>"
	}
}

func init() {
	var initEnv int32 = -1
	currentEnvironment = &initEnv
	environments := map[string]Environment{
		"prod":        PROD,
		"production":  PROD,
		"stage":       STAGE,
		"staging":     STAGE,
		"dev":         DEV,
		"development": DEV,
	}
	for envName, envType := range environments {
		if isEnvEquals(envName) {
			SetEnv(envType)
			internalLog.Infof("%s mode", envType)
			return
		}
	}
}

func isEnvEquals(rawEnv string) bool {
	grameEnv := strings.ToLower(os.Getenv("GRAMEWORK_ENV"))
	env := strings.ToLower(os.Getenv("ENV"))
	return strings.HasPrefix(grameEnv, rawEnv) || strings.HasPrefix(env, rawEnv)
}

// SetEnv sets gramework's environment
func SetEnv(e Environment) {
	if e != DEV && e != STAGE && e != PROD {
		internalLog.Warn("could not set unknown environment value, ignoring")
		return
	}
	if e != GetEnv() {
		internalLog.
			WithField("prevEnv", GetEnv()).
			WithField("newEnv", e).
			Warn("Setting a new environment")
	}
	if e == PROD {
		Logger.Level = log.InfoLevel
		enableDebug = false
	} else {
		enableDebug = true
		Logger.Level = log.DebugLevel
	}
	atomic.StoreInt32(currentEnvironment, int32(e))
}

// GetEnv returns current gramework's environment
func GetEnv() Environment {
	if currentEnvironment == nil {
		return DEV
	}
	return Environment(atomic.LoadInt32(currentEnvironment))
}

// FastHTTPLoggerAdapter Adapter for passing apex/log used as gramework Logger into fasthttp
type FastHTTPLoggerAdapter struct {
	apexLogger log.Interface
	fasthttp.Logger
}

// Logger handles default logger
var Logger = &log.Logger{
	Level:   log.InfoLevel,
	Handler: cli.New(os.Stdout),
}

// Errorf logs an error using default logger
func Errorf(msg string, v ...interface{}) {
	Logger.Errorf(msg, v...)
}

// NewFastHTTPLoggerAdapter create new *FastHTTPLoggerAdapter
func NewFastHTTPLoggerAdapter(logger *log.Interface) (fasthttplogger *FastHTTPLoggerAdapter) {
	fasthttplogger = &FastHTTPLoggerAdapter{
		apexLogger: *logger,
	}
	return fasthttplogger
}

// Printf show message only if set app.Logger.Level = apex/log.DebugLevel
func (l *FastHTTPLoggerAdapter) Printf(msg string, v ...interface{}) {
	l.apexLogger.Debugf(msg, v...)
}

var internalLog = func() *log.Entry {
	Logger.Level = log.DebugLevel
	if !enableDebug {
		Logger.Level = log.InfoLevel
	}

	return Logger.WithField("package", "gramework")
}()
