package gramework

// ListenAndServeAll serves HTTP and HTTPS automatically.
// HTTPS is served on :443.
// If it can't serve http or https, it logs an error and
// exit the server with app.Logger.Fatalf().
func (app *App) ListenAndServeAll(httpAddr ...string) {
	go func() {
		err := app.ListenAndServeAutoTLS(":443")
		app.Logger.Fatalf("can't serve tls: %s", err)
	}()
	err := app.ListenAndServe(httpAddr...)
	app.Logger.Fatalf("can't serve http: %s", err)
}

// ListenAndServeAllDev serves HTTP and HTTPS automatically
// with localhost HTTPS support via self-signed certs.
// HTTPS is served on :443.
// If it can't serve http or https, it logs an error and
// exit the server with app.Logger.Fatalf().
func (app *App) ListenAndServeAllDev(httpAddr ...string) {
	go func() {
		err := app.ListenAndServeAutoTLSDev(":443")
		app.Logger.Fatalf("can't serve tls: %s", err)
	}()
	err := app.ListenAndServe(httpAddr...)
	app.Logger.Fatalf("can't serve http: %s", err)
}
