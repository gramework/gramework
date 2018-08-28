// Copyright 2017-present Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package gramework

// Domain returns a domain router
func (app *App) Domain(domain string) *Router {
	app.domainListLock.Lock()
	if app.domains[domain] == nil {
		app.domains[domain] = &Router{
			router: newRouter(),
			app:    app,
		}
	}
	app.domainListLock.Unlock()
	return app.domains[domain]
}
