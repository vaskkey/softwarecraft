package main

import (
	"net/http"
)

// getRoot is a handler that shows index page
func (app *application) getRoot(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "home.tmpl.html", templateData{})
}
