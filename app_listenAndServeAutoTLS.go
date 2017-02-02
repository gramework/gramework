package gramework

import (
	"crypto/tls"
	"net"
	"strings"

	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/acme/autocert"
	"rsc.io/letsencrypt"
)

func (app *App) ListenAndServeAutoTLS(addr string, cachePath ...string) error {
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

	m := autocert.Manager{
		Prompt: autocert.AcceptTOS,
	}

	if letscache != "" {
		m.Cache = autocert.DirCache(letscache)
	}

	tlsConfig := &tls.Config{
		GetCertificate: m.GetCertificate,
	}
	tlsLn := tls.NewListener(ln, tlsConfig)

	err = fasthttp.Serve(tlsLn, app.router.Handler)
	if err != nil {
		app.Logger.Errorf("Can't serve: %s", err)
	}
	return err
}

func (app *App) ListenAndServeAutoTLSDev(addr string, cachePath ...string) error {
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
	m.Register("k@guava.by", func(string) bool { return true })

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

	err = fasthttp.Serve(tlsLn, app.router.Handler)
	if err != nil {
		app.Logger.Errorf("Can't serve: %s", err)
	}
	return err
}
