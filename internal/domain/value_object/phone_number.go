package value_object

import (
	"errors"
	"regexp"
)

type PhoneNumber struct {
	value string
}

var regex = regexp.MustCompile(`\D`)

func NewPhoneNumber(value string) (*PhoneNumber, error) {
	cleanedNumber := regex.ReplaceAllString(value, "")
	if cleanedNumber == "" {
		return nil, errors.New("value is empty")
	}
	numberLength := len(cleanedNumber)
	if numberLength != 10 && numberLength != 11 {
		return nil, errors.New("phone number length is invalid")
	}
	return &PhoneNumber{value: cleanedNumber}, nil
}

func (p *PhoneNumber) Value() *string {
	return &p.value
}
