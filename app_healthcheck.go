// Copyright 2017-present Kirill Danshin and Gramework contributors
// Copyright 2019-present Highload LTD (UK CN: 11893420)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package gramework

import (
	"runtime"
)

// HealthHandler serves info about memory usage
func (app *App) HealthHandler(ctx *Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	e := ctx.JSON(m)
	_ = e
}

// Health registers HealthHandler on /internal/health
func (app *App) Health() {
	app.GET("/internal/health", app.HealthHandler)
}
