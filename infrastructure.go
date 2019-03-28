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
	"time"

	"github.com/gramework/gramework/infrastructure"
)

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
