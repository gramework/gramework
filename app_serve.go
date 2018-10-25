// Copyright 2017-present Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package gramework

import (
	"net"
)

// Serve app on given listener
func (app *App) Serve(ln net.Listener) error {
	var err error
	srv := app.copyServer()
	app.runningServersMu.Lock()
	app.runningServers = append(app.runningServers, runningServerInfo{
		bind: ln.Addr().String(),
		srv:  srv,
	})
	app.runningServersMu.Unlock()
	if err = srv.Serve(ln); err != nil {
		app.internalLog.Errorf("ListenAndServe failed: %s", err)
	}

	return err
}
