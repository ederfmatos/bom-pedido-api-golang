package admin

import "bom-pedido-api/domain/value_object"

type Admin struct {
	Id         string
	Name       string
	Email      value_object.Email
	MerchantId string
}

func New(name, rawEmail, merchantId string) (*Admin, error) {
	email, err := value_object.NewEmail(rawEmail)
	if err != nil {
		return nil, err
	}
	return &Admin{
		Id:         value_object.NewID(),
		Name:       name,
		Email:      *email,
		MerchantId: merchantId,
	}, nil
}

func Restore(id, name, rawEmail, merchantId string) (*Admin, error) {
	email, err := value_object.NewEmail(rawEmail)
	if err != nil {
		return nil, err
	}
	return &Admin{
		Id:         id,
		Name:       name,
		Email:      *email,
		MerchantId: merchantId,
	}, nil
}

func (a *Admin) GetEmail() string {
	return a.Email.Value()
}
