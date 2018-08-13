// Copyright 2017 Kirill Danshin and Gramework contributors
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
)

// ListenAndServe HTTP on given addr.
// runs flag.Parse() if !flag.Parsed() to support
// --bind flag.
func (app *App) ListenAndServe(addr ...string) error {
	var bind string
	if len(addr) > 0 {
		bind = addr[0]
	} else {
		if !app.flagsRegistered {
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
		return errors.New("No bind address provided")
	}

	l := app.Logger.WithField("bind", bind)
	l.Info("Starting HTTP")

	if len(app.name) == 0 {
		app.name = "gramework/" + Version
	}

	var err error
	srv := app.copyServer()
	if err = srv.ListenAndServe(bind); err != nil {
		l.Errorf("ListenAndServe failed: %s", err)
	}

	return err
}
