package customer

import (
	"bom-pedido-api/internal/domain/value_object"
)

const (
	ACTIVE string = "ACTIVE"
)

type Customer struct {
	Id          string
	Name        string
	email       value_object.Email
	phoneNumber *value_object.PhoneNumber
	Status      string
	TenantId    string
}

func New(name, email, tenantId string) (*Customer, error) {
	newEmail, err := value_object.NewEmail(email)
	if err != nil {
		return nil, err
	}
	return &Customer{
		Id:       value_object.NewID(),
		Name:     name,
		email:    *newEmail,
		Status:   ACTIVE,
		TenantId: tenantId,
	}, nil
}

func Restore(id, name, email string, phoneNumber *string, status, tenantId string) (*Customer, error) {
	newEmail, err := value_object.NewEmail(email)
	if err != nil {
		return nil, err
	}
	var newPhoneNumber *value_object.PhoneNumber
	if phoneNumber != nil {
		newPhoneNumber, err = value_object.NewPhoneNumber(*phoneNumber)
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
		TenantId:    tenantId,
	}, nil
}

func (customer *Customer) GetPhoneNumber() *string {
	if customer.phoneNumber == nil {
		return nil
	}
	return customer.phoneNumber.Value()
}

func (customer *Customer) SetPhoneNumber(phoneNumber string) error {
	newPhoneNumber, err := value_object.NewPhoneNumber(phoneNumber)
	if err != nil {
		return err
	}
	customer.phoneNumber = newPhoneNumber
	return nil
}

func (customer *Customer) GetEmail() string {
	return customer.email.Value()
}
