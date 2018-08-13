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
	"crypto/tls"
	"math/rand"
	"net"
	"runtime"
	"strings"
	"time"

	"github.com/kirillDanshin/letsencrypt"
	"golang.org/x/crypto/acme/autocert"
)

func getDefaultTLSConfig() *tls.Config {
	return &tls.Config{
		DynamicRecordSizingDisabled: true,
		ClientSessionCache:          tls.NewLRUClientSessionCache(1024 * runtime.GOMAXPROCS(0)),
	}
}

func getCachePath(dev ...bool) string {
	p := "./tls-gramecache"
	if len(dev) > 0 && dev[0] {
		p += ".dev"
	}

	return p
}

// ListenAndServeAutoTLS serves TLS requests
func (app *App) ListenAndServeAutoTLS(addr string, cachePath ...string) error {
	var ln net.Listener
	var err error

	if len(app.TLSEmails) == 0 {
		return ErrTLSNoEmails
	}

	if portIdx := strings.IndexByte(addr, ':'); portIdx == -1 {
		addr += ":443"
	}

	ln, err = net.Listen("tcp4", addr)
	if err != nil {
		app.Logger.Errorf("Can't serve: %s", err)
		return err
	}

	letscache := getCachePath()
	if len(cachePath) > 0 {
		letscache = cachePath[0]
	}

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	m := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Email:  app.TLSEmails[r.Intn(len(app.TLSEmails))],
	}

	if letscache != "" {
		m.Cache = autocert.DirCache(letscache)
	}

	tlsConfig := getDefaultTLSConfig()
	tlsConfig.GetCertificate = func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		var cert *tls.Certificate
		cert, err = m.GetCertificate(hello)
		if err != nil {
			app.Logger.Errorf("can't get cert: %s", err)
		}

		return cert, err
	}

	tlsLn := tls.NewListener(ln, tlsConfig)

	l := app.Logger.WithField("bind", addr)
	l.Info("Starting HTTPS")

	if len(app.name) == 0 {
		app.name = "gramework/" + Version
	}

	srv := app.copyServer()
	if err = srv.Serve(tlsLn); err != nil {
		app.Logger.Errorf("Can't serve: %s", err)
	}

	return err
}

// ListenAndServeAutoTLSDev serves non-production grade TLS requests. Supports localhost.localdomain.
func (app *App) ListenAndServeAutoTLSDev(addr string, cachePath ...string) error {
	if len(app.TLSEmails) == 0 {
		return ErrTLSNoEmails
	}

	if portIdx := strings.IndexByte(addr, ':'); portIdx == -1 {
		addr += ":443"
	}

	var ln net.Listener
	var err error
	ln, err = net.Listen("tcp4", addr)
	if err != nil {
		app.Logger.Errorf("Can't serve: %s", err)
		return err
	}

	letscache := getCachePath(true)
	if len(cachePath) > 0 {
		letscache = cachePath[0]
	}

	var m letsencrypt.Manager
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	err = m.Register(
		app.TLSEmails[r.Intn(len(app.TLSEmails))],
		autocert.AcceptTOS,
	)
	if err != nil {
		return err
	}

	if letscache != "" {
		if err = m.CacheFile(letscache); err != nil {
			app.Logger.Errorf("Can't serve: %s", err)
			return err
		}
	}

	tlsConfig := getDefaultTLSConfig()
	tlsConfig.GetCertificate = func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		var cert *tls.Certificate
		cert, err = m.GetCertificate(hello)
		if err != nil {
			app.Logger.Errorf("can't get cert: %s", err)
		}

		return cert, err
	}

	tlsLn := tls.NewListener(ln, tlsConfig)

	l := app.Logger.WithField("bind", addr)
	l.Info("Starting HTTPS")

	srv := app.copyServer()
	if err = srv.Serve(tlsLn); err != nil {
		app.Logger.Errorf("Can't serve: %s", err)
	}

	return err
}
