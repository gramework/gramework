package gramework

// Shutdown gracefully shuts down application servers
func (app *App) Shutdown() (err error) {
	app.runningServersMu.Lock()
	// this is not a hot path, we can freely use defer here
	defer app.runningServersMu.Unlock()

	newRunningList := []runningServerInfo{}
	for _, info := range app.runningServers {
		app.Logger.WithField("bind", info.bind).Warn("shutting down server")
		err = info.srv.Shutdown()
		if err != nil {
			app.Logger.WithError(err).Error("could not shutdown server")
			newRunningList = append(newRunningList, info)
			continue
		}
	}

	app.runningServers = newRunningList

	if err == nil {
		app.Logger.Warn("application servers shutted down successfully")
		return
	}
	app.Logger.WithError(err).WithField("stillRunning", len(app.runningServers)).Warn("could not stop servers")
	return
}
