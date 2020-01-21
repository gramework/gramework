// Copyright 2017-present Kirill Danshin and Gramework contributors
// Copyright 2019-present Highload LTD (UK CN: 11893420)
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
	"fmt"
	"math/rand"
	"net"
	"runtime"
	"time"

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
	if len(app.TLSEmails) == 0 {
		return ErrTLSNoEmails
	}

	addr, err := normalizeTLSAddr(addr)
	if err != nil {
		app.internalLog.Errorf("Bad address %q: %s", addr, err)
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		app.internalLog.Errorf("Can't serve %q: %s", addr, err)
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
	tlsConfig.GetCertificate = func(hello *tls.ClientHelloInfo) (cert *tls.Certificate, err error) {
		if len(hello.ServerName) == 0 || hello.ServerName == localhost {
			hello.ServerName = localhost
			cert, err = selfSignedCertificate(hello)
		} else {
			cert, err = m.GetCertificate(hello)
		}

		if err != nil {
			app.internalLog.Errorf("Can't get cert for %q: %s", hello.ServerName, err)
		}

		return cert, err
	}

	tlsLn := tls.NewListener(ln, tlsConfig)
	checks()

	l := app.internalLog.WithField("bind", addr)
	l.Info("Starting HTTPS")

	srv := app.copyServer()
	app.runningServersMu.Lock()
	app.runningServers = append(app.runningServers, runningServerInfo{
		bind: addr,
		srv:  srv,
	})
	app.runningServersMu.Unlock()
	if err = srv.Serve(tlsLn); err != nil {
		app.internalLog.Errorf("Can't serve: %s", err)
	}

	return err
}

// ListenAndServeAutoTLSDev serves non-production grade TLS requests. Supports localhost.localdomain.
// Deprecated: use ListenAndServeAutoTLS() instead
func (app *App) ListenAndServeAutoTLSDev(addr string, cachePath ...string) error {
	return app.ListenAndServeAutoTLS(addr, cachePath...)
}

func normalizeTLSAddr(addr string) (string, error) {
	host, port, err := net.SplitHostPort(addr)

	if err != nil && net.ParseIP(host) == nil {
		return addr, fmt.Errorf("invalid address %q: %s", addr, err)
	}

	if len(port) == 0 {
		addr = net.JoinHostPort(host, "443")
	}

	return addr, nil
}
