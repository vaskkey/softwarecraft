package main

import "net/http"

// routes maps routes to the handlers, maps the static files to enpoints.
// and returns the ServeMux.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("GET /register", app.getRegister)
	mux.HandleFunc("POST /register", app.postRegister)

	mux.HandleFunc("GET /login", app.getLogin)
	mux.HandleFunc("POST /login", app.postLogin)

	// Root page
	mux.HandleFunc("GET /{$}", app.getRoot)

	// Serving static files
	fs := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fs))

	return app.recoverPanic(app.logRequest(commonHeaders(mux)))
}
