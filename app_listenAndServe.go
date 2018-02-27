package gramework

import (
	"errors"
	"flag"

	"github.com/valyala/fasthttp"
)

// ListenAndServe HTTP on given addr.
// runs flag.Parse() if !flag.Parsed() to support
// --bind flag.
func (app *App) ListenAndServe(addr ...string) error {
	var bind string
	if len(addr) > 0 {
		bind = addr[0]
	} else {
		if !app.flagsRegistered {
			app.RegFlags()
		}
	}
	if !flag.Parsed() {
		flag.Parse()
	}
	if app.Flags.values != nil {
		if bindFlag, ok := app.Flags.values["bind"]; ok {
			bind = *bindFlag.Value
		}
	}
	if bind == "" {
		return errors.New("No bind address provided")
	}
	l := app.Logger.WithField("bind", bind)

	l.Info("Starting HTTP")

	if len(app.name) == 0 {
		app.name = "gramework/" + Version
	}

	s := fasthttp.Server{
		Handler: app.handler(),
		Logger:  NewFastHTTPLoggerAdapter(&app.Logger),
		Name:    app.name,
	}
	err := s.ListenAndServe(bind)
	l.Errorf("ListenAndServe failed: %s", err)
	return err
}
