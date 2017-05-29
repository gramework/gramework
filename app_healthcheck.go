package gramework

import (
	"runtime"
)

// HealthHandler serves info about memory usage
func (app *App) HealthHandler(ctx *Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	ctx.JSON(m)
}

// Health registers HealthHandler on /internal/health
func (app *App) Health() {
	app.GET("/internal/health", app.HealthHandler)
}
