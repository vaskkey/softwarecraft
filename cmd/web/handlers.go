package main

import (
	"net/http"

	"github.com/vaskkey/softwarecraft/internal/config"
	"github.com/vaskkey/softwarecraft/internal/models"
)

// getRoot renders main page
func (app *application) getRoot(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "home.tmpl.html", templateData{})
}

// getRegister renders login page
func (app *application) getRegister(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "register.tmpl.html", templateData{})
}

// postRegister saves user to DB
func (app *application) postRegister(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	uParams := models.NewUserParams(&r.PostForm)

	if ok, errs := uParams.Validate(); !ok {
		app.render(w, r, http.StatusUnprocessableEntity, "register.tmpl.html", templateData{
			Errs:       errs,
			UserParams: *uParams,
		})
		return
	}

	user, err := uParams.GetUser()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.users.Insert(user)
	if err != nil {
		switch err {
		case config.ErrDuplicateEmail:
			errs := make(config.ValidationErrors)
			errs["email"] = "User with this email already exists."
			app.render(w, r, http.StatusUnprocessableEntity, "register.tmpl.html", templateData{
				Errs:       errs,
				UserParams: *uParams,
			})
			return
		default:
			app.serverError(w, r, err)
			return
		}
	}

	app.render(w, r, http.StatusOK, "register.tmpl.html", templateData{})
}
