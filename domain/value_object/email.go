package value_object

import (
	"bom-pedido-api/domain/errors"
	"regexp"
)

type Email struct {
	value string
}

func NewEmail(email string) (*Email, error) {
	if len(email) == 0 {
		return nil, errors.EmailIsEmptyError
	}
	if len(email) > 255 {
		return nil, errors.EmailIsTooLongError
	}
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return nil, errors.InvalidEmailError
	}
	return &Email{value: email}, nil
}

func (e Email) Value() *string {
	return &e.value
}
