package main

import (
	"context"
	"net/http"

	"github.com/vaskkey/softwarecraft/internal/models"
)

type contextKey string

const userContextKey = contextKey("user")

// contextSetUser set user in request context
func (app *application) contextSetUser(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// contextGetUser get user stored in request context
func (app *application) contextGetUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok {
		return nil
	}

	return user
}
