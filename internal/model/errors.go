package model

import "errors"

// Validation errors
var (
	ErrMissingName          = errors.New("name is required")
	ErrMissingApplicationID = errors.New("application_id is required")

	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrInvalidSignature = errors.New("invalid token signature")
	ErrMissingClaims    = errors.New("missing required claims")

	ErrPasswordTooLong = errors.New("password exceeds maximum length of 32 bytes")
	ErrPasswordEmpty   = errors.New("password cannot be empty")
)
