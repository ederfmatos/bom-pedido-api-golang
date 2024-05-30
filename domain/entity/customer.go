package entity

import "bom-pedido-api/domain/value_object"

const (
	ACTIVE   string = "ACTIVE"
	INACTIVE string = "INACTIVE"
)

type Customer struct {
	Id          string
	Name        string
	Email       value_object.Email
	PhoneNumber *value_object.PhoneNumber
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
		Email:  *newEmail,
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
		Email:       *newEmail,
		PhoneNumber: newPhoneNumber,
		Status:      status,
	}, nil
}

func (customer Customer) isActive() bool {
	return customer.Status == ACTIVE
}

func (customer Customer) isInactive() bool {
	return customer.Status == INACTIVE
}
