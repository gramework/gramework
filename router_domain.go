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
