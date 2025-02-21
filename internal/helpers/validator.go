package helpers

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type ValidationErrors = map[string]string

type Validator struct {
	Errors ValidationErrors
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddFieldError add the error to the errors hashmap
func (v *Validator) AddFieldError(key, message string) {
	if v.Errors == nil {
		v.Errors = make(map[string]string)
	}

	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// CheckField if there is an error - add the message to hashmap under the given key
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// MaxLength returns true if number of chars in value is less or equal than length
func (v *Validator) MaxLength(value string, length int) bool {
	return utf8.RuneCountInString(value) <= length
}

// MinLength returns true if number of chars in value is more or equal than length
func (v *Validator) MinLength(value string, length int) bool {
	return utf8.RuneCountInString(value) >= length
}

// LengthBetween checks if number of characters is between the bounds
func (v *Validator) LengthBetween(value string, from, to int) bool {
	return v.MinLength(value, from) && v.MaxLength(value, to)
}

// IsNotBlank checks if stripped string is not empty
func (v *Validator) IsNotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}
