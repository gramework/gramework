// Package xhostname middleware provides `X-Hostname` header in each request
// and useful when using scalable container platform to see
// which host sent you current response.
//
// To use this package, just call `xhostname.Setup`:
//
//		app := gramework.New()
//		xhostname.Setup(app)
package xhostname

import (
	"os"

	"github.com/gramework/gramework"
)

const (
	// HeaderKey is the header name which we use
	HeaderKey = "X-Hostname"
)

var hostname string

func init() {
	// kubernetes provides a HOSTNAME in pods
	if h := os.Getenv("HOSTNAME"); len(h) > 0 {
		hostname = h
		gramework.Logger.WithField("hostname", h).Info("using HOSTNAME env as a hostname")
		return
	}
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		gramework.Logger.WithError(err).Error("could not get hostname")
		return
	}
	gramework.Logger.WithField("hostname", hostname).Info("using os.Hostname env as a hostname")
}

// Setup registers middleware in the provided app
func Setup(app *gramework.App) {
	err := app.UseAfterRequest(serveXHost)
	if err != nil {
		app.Logger.WithError(err).WithField("package", "mw/xhostname").Error("could not register middleware")
	}
}

func serveXHost(ctx *gramework.Context) {
	ctx.Response.Header.Add(HeaderKey, hostname)
}
