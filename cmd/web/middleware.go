package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/vaskkey/softwarecraft/internal/helpers"
)

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CSP https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP
		// w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com unpkg.com; font-src fonts.gstatic.com; connect-src unpkg.com/boxicons@2.1.4/svg/")
		// Same origin requests https://developer.mozilla.org/en-US/docs/Web/Security/Same-origin_policy
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		// Mime sniffing https://security.stackexchange.com/questions/7506/using-file-extension-and-mime-type-as-output-by-file-i-b-combination-to-dete/7531#7531
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// Clickjacking protection https://developer.mozilla.org/en-US/docs/Web/Security/Types_of_attacks#click-jacking
		w.Header().Set("X-Frame-Options", "deny")
		// Recommended to disable since we're enabling CSP options above https://owasp.org/www-project-secure-headers/#x-xss-protection
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("received request", slog.String("ip", r.RemoteAddr), slog.String("proto", r.Proto), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Recover from panic and return 500
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) authorizedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ID, ok := app.sessionManager.Get(r.Context(), "userID").(int64)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := app.users.GetByID(ID)
		if err != nil {
			switch {
			case errors.Is(err, helpers.ErrNoRecords):
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			default:
				app.serverError(w, r, err)
				return
			}
		}

		r = app.contextSetUser(r, &user)

		next.ServeHTTP(w, r)
	})
}
