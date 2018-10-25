// Copyright 2017-present Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package gramework

import "github.com/valyala/fasthttp"

// OptAppName sets app.name and app.serverBase.Name
func OptAppName(n string) func(*App) {
	return func(a *App) {
		if a != nil && a.serverBase != nil {
			a.name = n
			a.serverBase.Name = n
		}
	}
}

// OptUseServer sets fasthttp.Server instance to use
func OptUseServer(s *fasthttp.Server) func(*App) {
	return func(a *App) {
		if a != nil && s != nil {
			a.serverBase = s
			a.serverBase.Handler = a.handler()
		}
	}
}

// OptMaxRequestBodySize sets new MaxRequestBodySize in the server used at the execution time.
// All OptUseServer will overwrite this setting 'case OptUseServer replaces the whole server instance
// with a new one.
func OptMaxRequestBodySize(new int) func(*App) {
	return func(a *App) {
		if a != nil && a.serverBase != nil {
			a.serverBase.MaxRequestBodySize = new
		}
	}
}
