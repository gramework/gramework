package gramework

import (
	"io/ioutil"
	"log"
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
		Logger:  fasthttp.Logger(log.New(ioutil.Discard, "", log.LstdFlags)),
		Name:    app.name,
	}
	err := s.Serve(ln)
	app.Logger.Errorf("ListenAndServe failed: %s", err)
	return err
}
