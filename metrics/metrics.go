package metrics

import (
	"time"

	"github.com/gramework/gramework"
)

// Middleware handles metrics data
type Middleware struct {
}

// Register the middlewares
func Register(app *gramework.App) {
	m := Middleware{}
	app.UsePre(m.startReq)
	app.UseAfterRequest(m.endReq)
}

func (m *Middleware) startReq(ctx *gramework.Context) {
	ctx.SetUserValue("gramework.metrics.startTime", time.Now())
}

func (m *Middleware) endReq(ctx *gramework.Context) {
	startTime := ctx.UserValue("gramework.metrics.startTime").(time.Time)
	ctx.Logger.Infof("Served in %s", time.Now().Sub(startTime))
}
