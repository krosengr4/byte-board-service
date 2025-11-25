package model

import "errors"

// Validation errors
var (
	ErrMissingName          = errors.New("name is required")
	ErrMissingApplicationID = errors.New("application_id is required")
)
