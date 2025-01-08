package customer

import (
	"bom-pedido-api/internal/domain/value_object"
)

const (
	ACTIVE string = "ACTIVE"
)

type Customer struct {
	Id          string                   `bson:"id"`
	Name        string                   `bson:"name"`
	Email       value_object.Email       `bson:"email"`
	PhoneNumber value_object.PhoneNumber `bson:"phoneNumber"`
	Status      string                   `bson:"status"`
	TenantId    string                   `bson:"tenantId"`
}

func New(name, email, tenantId string) (*Customer, error) {
	newEmail, err := value_object.NewEmail(email)
	if err != nil {
		return nil, err
	}
	return &Customer{
		Id:       value_object.NewID(),
		Name:     name,
		Email:    *newEmail,
		Status:   ACTIVE,
		TenantId: tenantId,
	}, nil
}

func (customer *Customer) GetPhoneNumber() *string {
	return customer.PhoneNumber.Value()
}

func (customer *Customer) SetPhoneNumber(phoneNumber string) error {
	newPhoneNumber, err := value_object.NewPhoneNumber(phoneNumber)
	if err != nil {
		return err
	}
	customer.PhoneNumber = *newPhoneNumber
	return nil
}

func (customer *Customer) GetEmail() string {
	return customer.Email.Value()
}
