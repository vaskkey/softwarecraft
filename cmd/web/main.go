package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/vaskkey/softwarecraft/internal/models"
)

type application struct {
	// Utils
	logger         *slog.Logger
	templateCache  templateCache
	sessionManager *scs.SessionManager

	// Models
	users *models.UserModel
}

func main() {
	// ENV
	addr := flag.String("addr", ":8080", "HTTP network address")
	dsn := flag.String("dsn", "postgres://postgres:postgres@localhost:5432/softwarecraft?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// DB connection
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// HTML template setup
	tc, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Session setup
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		logger:         logger,
		templateCache:  tc,
		users:          models.NewUserModel(db),
		sessionManager: sessionManager,
	}

	logger.Info("Starting server on", "port", *addr)
	err = http.ListenAndServe(*addr, app.routes())

	logger.Error(err.Error())
	os.Exit(1)
}
