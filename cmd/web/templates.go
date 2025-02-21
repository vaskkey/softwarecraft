package main

import (
	"html/template"
	"path/filepath"

	"github.com/vaskkey/softwarecraft/internal/helpers"
	"github.com/vaskkey/softwarecraft/internal/models"
)

type templateCache map[string]*template.Template

type templateData struct {
	UserParams models.RegisterUser
	Errs       helpers.ValidationErrors
}

// newTemplateCache parses all files in templates folder and prepares them to be rendered
func newTemplateCache() (templateCache, error) {
	cache := templateCache{}

	// Get templates for base app
	err := handleTemplateDir("./ui/html/pages/*.tmpl.html", "./ui/html/base.tmpl.html", &cache)
	if err != nil {
		return nil, err
	}

	// Get templates for unauthenticated parts of the app
	err = handleTemplateDir("./ui/html/auth-pages/*.tmpl.html", "./ui/html/auth_base.tmpl.html", &cache)
	if err != nil {
		return nil, err
	}

	return cache, nil
}

// handleTemplateDir reusable function that allows us to parse different directories
func handleTemplateDir(globPath, basePath string, cache *templateCache) error {
	pages, err := filepath.Glob(globPath)
	if err != nil {
		return err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles(basePath)
		if err != nil {
			return err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return err
		}

		(*cache)[name] = ts
	}

	return nil
}
