package gramework

import (
	"net"

	"github.com/valyala/fasthttp"
)

// Serve app on given listener
func (app *App) Serve(ln net.Listener) error {
	if len(app.name) == 0 {
		app.name = "gramework/" + Version
	}

	s := fasthttp.Server{
		Handler: app.handler(),
		Logger:  NewFastHTTPLoggerAdapter(&app.Logger),
		Name:    app.name,
	}

	var err error
	if err = s.Serve(ln); err != nil {
		app.Logger.Errorf("ListenAndServe failed: %s", err)
	}

	return err
}
