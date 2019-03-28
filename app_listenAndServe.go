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
	"errors"
	"flag"
	"os"
	"strings"
)

// ListenAndServe HTTP on given addr.
// runs flag.Parse() if !flag.Parsed() to support
// --bind flag.
func (app *App) ListenAndServe(addr ...string) error {
	var bind string
	if len(addr) > 0 {
		bind = addr[0]
	} else {
		bind = os.Getenv("PORT")
		if len(bind) > 0 && !strings.Contains(bind, ":") {
			bind = ":" + bind
		}
		if bind == "" && !app.flagsRegistered {
			app.RegFlags()
		}
	}

	if !flag.Parsed() && !flagsDisabled {
		flag.Parse()
	}

	if app.Flags.values != nil {
		if bindFlag, ok := app.Flags.values["bind"]; ok {
			bind = *bindFlag.Value
		}
	}

	if bind == "" {
		return errors.New("no bind address provided")
	}

	l := app.internalLog.WithField("bind", bind)
	l.Info("Starting HTTP")

	var err error
	srv := app.copyServer()
	app.runningServersMu.Lock()
	app.runningServers = append(app.runningServers, runningServerInfo{
		bind: bind,
		srv:  srv,
	})
	app.runningServersMu.Unlock()
	if err = srv.ListenAndServe(bind); err != nil {
		l.Errorf("ListenAndServe failed: %s", err)
	}

	return err
}
