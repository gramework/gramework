package gramework

func (app *App) ListenAndServeAll(addr ...string) {
	go func() {
		err := app.ListenAndServeAutoTLS(":443")
		app.Logger.Fatalf("can't serve tls: %s", err)
	}()
	err := app.ListenAndServe(addr...)
	app.Logger.Fatalf("can't serve http: %s", err)
}

func (app *App) ListenAndServeAllDev(addr ...string) {
	go func() {
		err := app.ListenAndServeAutoTLSDev(":443")
		app.Logger.Fatalf("can't serve tls: %s", err)
	}()
	err := app.ListenAndServe(addr...)
	app.Logger.Fatalf("can't serve http: %s", err)
}
