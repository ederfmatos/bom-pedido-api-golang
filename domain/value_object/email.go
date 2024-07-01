package value_object

import (
	"bom-pedido-api/domain/errors"
	"regexp"
)

var (
	EmailIsEmptyError   = errors.New("email is empty")
	EmailIsTooLongError = errors.New("email is too long")
	InvalidEmailError   = errors.New("invalid email")
)

type Email struct {
	value string
}

func NewEmail(email string) (*Email, error) {
	if len(email) == 0 {
		return nil, EmailIsEmptyError
	}
	if len(email) > 255 {
		return nil, EmailIsTooLongError
	}
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return nil, InvalidEmailError
	}
	return &Email{value: email}, nil
}

func (e Email) Value() *string {
	return &e.value
}
