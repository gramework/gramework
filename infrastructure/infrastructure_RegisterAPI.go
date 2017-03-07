package infrastructure

import (
	"time"

	"github.com/gramework/gramework"
)

// RegisterAPI in the app
func (i *Infrastructure) RegisterAPI(app *gramework.App) {
	app.GET("/infrastructure", func(ctx *gramework.Context) {
		i.Lock.RLock()
		i.CurrentTimestamp = time.Now().UnixNano()
		ctx.JSON(i)
		i.Lock.RUnlock()
	})
}
