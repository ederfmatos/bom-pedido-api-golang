package errors

import "errors"

var (
	EmailIsEmptyError   = errors.New("email is empty")
	EmailIsTooLongError = errors.New("email is too long")
	InvalidEmailError   = errors.New("invalid email")
)
