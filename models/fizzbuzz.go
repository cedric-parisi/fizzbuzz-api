package models

import (
	"fmt"

	"github.com/cedric-parisi/fizzbuzz-api/internal/errors"
)

const (
	maxLimit = 1000
)

// Fizzbuzz holds a fizzbuzz query data.
type Fizzbuzz struct {
	Int1  int    `json:"int1"`
	Int2  int    `json:"int2"`
	Limit int    `json:"limit"`
	Str1  string `json:"str1"`
	Str2  string `json:"str2"`
}

// Validate validates required fields for fizzbuzz.
func (f Fizzbuzz) Validate() []errors.FieldError {
	var errs []errors.FieldError

	if f.Int1 <= 0 {
		errs = append(errs, errors.FieldError{
			Name:    "int1",
			Message: "int1 must be greater than 0",
		})
	}
	if f.Int2 <= 0 {
		errs = append(errs, errors.FieldError{
			Name:    "int2",
			Message: "int2 must be greater than 0",
		})
	}
	if f.Int2 <= f.Int1 {
		errs = append(errs, errors.FieldError{
			Name:    "int2",
			Message: "int2 must be greater than int1",
		})
	}

	if f.Limit <= 0 {
		errs = append(errs, errors.FieldError{
			Name:    "limit",
			Message: "limit must be greater than 0",
		})
	}

	if f.Limit > maxLimit {
		errs = append(errs, errors.FieldError{
			Name:    "limit",
			Message: fmt.Sprintf("limit must be smaller than %d", maxLimit),
		})
	}

	if f.Str1 == "" {
		errs = append(errs, errors.FieldError{
			Name:    "str1",
			Message: "str1 must not be empty",
		})
	}

	if f.Str2 == "" {
		errs = append(errs, errors.FieldError{
			Name:    "str2",
			Message: "str2 must not be empty",
		})
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}
