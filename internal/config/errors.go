package config

import "errors"

var ErrNoRecords = errors.New("models: no record found")
var ErrDuplicateEmail = errors.New("models: user with this email already exists")

type ValidationErrors = map[string]string
