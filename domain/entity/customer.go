package entity

import (
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
)

const (
	ACTIVE   string = "ACTIVE"
	INACTIVE string = "INACTIVE"
)

var (
	CustomerNotFoundError = errors.New("customer not found")
)

type Customer struct {
	Id          string
	Name        string
	email       value_object.Email
	phoneNumber *value_object.PhoneNumber
	Status      string
}

func NewCustomer(name, email string) (*Customer, error) {
	newEmail, err := value_object.NewEmail(email)
	if err != nil {
		return nil, err
	}
	return &Customer{
		Id:     value_object.NewID(),
		Name:   name,
		email:  *newEmail,
		Status: ACTIVE,
	}, nil
}

func RestoreCustomer(id, name, email, phoneNumber, status string) (*Customer, error) {
	newEmail, err := value_object.NewEmail(email)
	if err != nil {
		return nil, err
	}
	var newPhoneNumber *value_object.PhoneNumber
	if phoneNumber != "" {
		newPhoneNumber, err = value_object.NewPhoneNumber(phoneNumber)
		if err != nil {
			return nil, err
		}
	}
	return &Customer{
		Id:          id,
		Name:        name,
		email:       *newEmail,
		phoneNumber: newPhoneNumber,
		Status:      status,
	}, nil
}

func (customer Customer) isActive() bool {
	return customer.Status == ACTIVE
}

func (customer Customer) isInactive() bool {
	return customer.Status == INACTIVE
}

func (customer Customer) GetPhoneNumber() *string {
	if customer.phoneNumber == nil {
		return nil
	}
	return customer.phoneNumber.Value()
}

func (customer Customer) SetPhoneNumber(phoneNumber string) error {
	newPhoneNumber, err := value_object.NewPhoneNumber(phoneNumber)
	if err != nil {
		return err
	}
	customer.phoneNumber = newPhoneNumber
	return nil
}

func (customer Customer) GetEmail() *string {
	return customer.email.Value()
}
