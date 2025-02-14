package auth

import (
	"errors"
	"fmt"
)

// AuthError is used to pass an error during the request through the
// application with auth specific context.
type AuthError struct {
	Msg string `json:"error"`
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (ae *AuthError) Error() string {
	return ae.Msg
}

// NewAuthError creates an AuthError for the provided message.
func NewAuthError(format string, args ...any) error {
	return &AuthError{
		Msg: fmt.Sprintf(format, args...),
	}
}

// IsAuthError checks if an error of type AuthError exists.
func IsAuthError(err error) bool {
	var ae *AuthError
	return errors.As(err, &ae)
}