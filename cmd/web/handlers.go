package main

import (
	"errors"
	"net/http"

	"github.com/vaskkey/softwarecraft/internal/helpers"
	"github.com/vaskkey/softwarecraft/internal/models"
)

// getRoot renders main page
func (app *application) getRoot(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "home.tmpl.html", templateData{})
}

// getRegister renders register page
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

	uParams := models.NewRegisterParams(&r.PostForm)

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
		case helpers.ErrDuplicateEmail:
			errs := make(helpers.ValidationErrors)
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

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// getLogin renders login page
func (app *application) getLogin(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "login.tmpl.html", templateData{})
}

// postLogin renders login page
func (app *application) postLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	errs := make(helpers.ValidationErrors)
	errs["login"] = "Unable to log in."

	uParams := models.NewLoginParams(&r.PostForm)
	if ok, _ := uParams.Validate(); !ok {
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", templateData{
			Errs: errs,
		})
		return
	}

	user, err := app.users.GetByEmail(uParams.Email)
	if err != nil {
		switch {
		case errors.Is(err, helpers.ErrNoRecords):
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", templateData{
				Errs: errs,
			})
			return
		default:
			app.serverError(w, r, err)
			return
		}
	}

	if match := user.Password.Compare(uParams.Password); !match {
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", templateData{
			Errs: errs,
		})
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
