package value_object

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
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

func (p *PhoneNumber) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if p == nil {
		return bson.MarshalValue("")
	}
	return bson.MarshalValue(p.value)
}

func (p *PhoneNumber) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	valueReader := bsonrw.NewBSONValueReader(t, data)
	value, err := valueReader.ReadString()
	if err != nil {
		return err
	}
	if value == "" {
		return nil
	}
	phoneNumber, err := NewPhoneNumber(value)
	if err != nil {
		return err
	}
	*p = *phoneNumber
	return nil
}
