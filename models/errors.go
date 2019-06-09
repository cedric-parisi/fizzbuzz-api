package models

import (
	"fmt"
)

var (
	// ErrNotFound is the error returned when the resource is not found.
	ErrNotFound = fmt.Errorf("not found")
)
