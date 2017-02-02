package gramework

import (
	"github.com/valyala/fasthttp"
)

// ListenAndServe on given addr
func (app *App) ListenAndServe(addr string) error {
	return fasthttp.ListenAndServe(addr, app.router.Handler)
}
