package main

import "net/http"

// routes maps routes to the handlers, maps the static files to enpoints.
// and returns the ServeMux.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Content Books
	mux.HandleFunc("/{$}", app.getRoot)

	return app.recoverPanic(app.logRequest(commonHeaders(mux)))
}
