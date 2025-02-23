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

	app.sessionManager.Put(r.Context(), "toast", "Successfully created new user.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// getLogin renders login page
func (app *application) getLogin(w http.ResponseWriter, r *http.Request) {
	toast := app.sessionManager.PopString(r.Context(), "toast")
	app.render(w, r, http.StatusOK, "login.tmpl.html", templateData{ToastMessage: toast})
}

// postLogin authenticates the user and puts their id in the session
func (app *application) postLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	td := templateData{
		ToastMessage: "Unable to log in.",
	}

	uParams := models.NewLoginParams(&r.PostForm)
	if ok, _ := uParams.Validate(); !ok {
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", td)
		return
	}

	user, err := app.users.GetByEmail(uParams.Email)
	if err != nil {
		switch {
		case errors.Is(err, helpers.ErrNoRecords):
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", td)
			return
		default:
			app.serverError(w, r, err)
			return
		}
	}

	if match := user.Password.Compare(uParams.Password); !match {
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", td)
		return
	}

	if err := app.sessionManager.RenewToken(r.Context()); err != nil {
		app.serverError(w, r, err)
	}

	app.sessionManager.Put(r.Context(), "userID", user.ID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// postLogout authenticates the user and puts their id in the session
func (app *application) postLogout(w http.ResponseWriter, r *http.Request) {
	if err := app.sessionManager.RenewToken(r.Context()); err != nil {
		app.serverError(w, r, err)
	}

	app.sessionManager.Remove(r.Context(), "userID")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
