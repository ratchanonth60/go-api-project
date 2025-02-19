package request

import "errors"

var (
	ErrEmailRequired    = errors.New("email is required")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrUsernameRequired = errors.New("username is required")
	ErrPasswordRequired = errors.New("password is required") // Add password validation
)
