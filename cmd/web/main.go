package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/vaskkey/softwarecraft/internal/models"
)

type application struct {
	logger        *slog.Logger
	templateCache templateCache
	users         *models.UserModel
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP network address")
	dsn := flag.String("dsn", "postgres://postgres:postgres@localhost:5432/softwarecraft?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	tc, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger:        logger,
		templateCache: tc,
		users:         models.NewUserModel(db),
	}

	logger.Info("Starting server on", "port", *addr)
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
