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

import "fmt"

// ListenAndServeAll serves HTTP and HTTPS automatically.
// HTTPS is served on :443.
// If it can't serve http or https, it logs an error and
// exit the server with app.Logger.Fatalf().
func (app *App) ListenAndServeAll(httpAddr ...string) {
	go func() {
		bindAddr := ":443"
		if app.TLSPort != 0 {
			bindAddr = fmt.Sprintf(":%d", app.TLSPort)
		}
		err := app.ListenAndServeAutoTLS(bindAddr)
		app.internalLog.Fatalf("can't serve tls: %s", err)
	}()

	if err := app.ListenAndServe(httpAddr...); err != nil {
		app.internalLog.Fatalf("can't serve http: %s", err)
	}
}

// ListenAndServeAllDev serves HTTP and HTTPS automatically
// with localhost HTTPS support via self-signed certs.
// HTTPS is served on :443.
// If it can't serve http or https, it logs an error and
// exit the server with app.Logger.Fatalf().
// Deprecated: Use ListenAndServeAll() instead
func (app *App) ListenAndServeAllDev(httpAddr ...string) {
	app.ListenAndServeAll(httpAddr...)
}
