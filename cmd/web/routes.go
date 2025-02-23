package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// routes maps routes to the handlers, maps the static files to enpoints.
// and returns the ServeMux.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Set up session middleware
	withSession := alice.New(app.sessionManager.LoadAndSave)

	// Auth routes
	mux.Handle("GET /register", withSession.ThenFunc(app.getRegister))
	mux.Handle("POST /register", withSession.ThenFunc(app.postRegister))

	mux.Handle("GET /login", withSession.ThenFunc(app.getLogin))
	mux.Handle("POST /login", withSession.ThenFunc(app.postLogin))
	mux.Handle("POST /logout", withSession.ThenFunc(app.postLogout))

	// Root page
	mux.HandleFunc("GET /{$}", app.getRoot)

	// Serving static files
	fs := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fs))

	main := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return main.Then(mux)
}
