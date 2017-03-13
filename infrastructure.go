package gramework

import (
	"time"

	"github.com/gramework/gramework/infrastructure"
)

var infrastructureServiceRegistrationErr = map[string]string{
	"error": "can't parse the query",
}

// ServeInfrastructure serves Infrastructure info
// It's an integration of our module
func (app *App) ServeInfrastructure(i *infrastructure.Infrastructure) {
	app.GET("/infrastructure", func(ctx *Context) {
		i.Lock.RLock()
		i.CurrentTimestamp = time.Now().UnixNano()
		ctx.CORS()
		ctx.JSON(i)
		i.Lock.RUnlock()
	})
	app.POST("/infrastructure/register/service", func(ctx *Context) {
		s := infrastructure.Service{
			Addresses: make([]infrastructure.Address, 0),
		}
		_, err := ctx.UnJSONBytes(ctx.PostBody(), &s)
		if err != nil {
			ctx.JSONError(err.Error())
			return
		}
		i.MergeService(s.Name, s)
		ctx.JSON(s)
	})
}
