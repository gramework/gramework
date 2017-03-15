package gramework

import (
	"crypto/tls"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/kirillDanshin/letsencrypt"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/acme/autocert"
)

// ListenAndServeAutoTLS serves TLS requests
func (app *App) ListenAndServeAutoTLS(addr string, cachePath ...string) error {
	if len(app.TLSEmails) == 0 {
		return ErrTLSNoEmails
	}
	if portIdx := strings.IndexByte(addr, ':'); portIdx == -1 {
		addr += ":443"
	}

	ln, err := net.Listen("tcp4", addr)
	if err != nil {
		app.Logger.Errorf("Can't serve: %s", err)
		return err
	}

	letscache := "./letscache"
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

	tlsConfig := &tls.Config{
		GetCertificate: m.GetCertificate,
	}
	tlsLn := tls.NewListener(ln, tlsConfig)

	l := app.Logger.WithField("bind", addr)
	l.Info("Starting HTTPS")
	err = fasthttp.Serve(tlsLn, app.handler())
	if err != nil {
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

	ln, err := net.Listen("tcp4", addr)
	if err != nil {
		app.Logger.Errorf("Can't serve: %s", err)
		return err
	}

	letscache := "./letscache.dev"
	if len(cachePath) > 0 {
		letscache = cachePath[0]
	}

	var m letsencrypt.Manager
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	m.Register(app.TLSEmails[r.Intn(len(app.TLSEmails))], func(string) bool { return true })

	if letscache != "" {
		if err = m.CacheFile(letscache); err != nil {
			app.Logger.Errorf("Can't serve: %s", err)
			return err
		}
	}

	tlsConfig := &tls.Config{
		GetCertificate: m.GetCertificate,
	}
	tlsLn := tls.NewListener(ln, tlsConfig)

	l := app.Logger.WithField("bind", addr)
	l.Info("Starting HTTPS")
	err = fasthttp.Serve(tlsLn, app.handler())
	if err != nil {
		app.Logger.Errorf("Can't serve: %s", err)
	}
	return err
}
