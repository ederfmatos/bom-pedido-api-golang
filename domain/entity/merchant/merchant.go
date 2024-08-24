package merchant

import (
	"bom-pedido-api/domain/value_object"
)

const (
	ACTIVE   Status = "ACTIVE"
	INACTIVE Status = "INACTIVE"
)

type (
	Status string

	Merchant struct {
		Id          string
		Name        string
		Email       value_object.Email
		PhoneNumber value_object.PhoneNumber
		Domain      string
		TenantId    string
		Status      Status
	}
)

func New(name, rawEmail, rawPhoneNumber, domain string) (*Merchant, error) {
	email, err := value_object.NewEmail(rawEmail)
	if err != nil {
		return nil, err
	}
	phoneNumber, err := value_object.NewPhoneNumber(rawPhoneNumber)
	if err != nil {
		return nil, err
	}
	return &Merchant{
		Id:          value_object.NewID(),
		Name:        name,
		Email:       *email,
		PhoneNumber: *phoneNumber,
		Domain:      domain,
		TenantId:    value_object.NewTenantId(),
		Status:      ACTIVE,
	}, nil
}

func Restore(id, name, rawEmail, rawPhoneNumber, status, domain, tenantId string) (*Merchant, error) {
	email, err := value_object.NewEmail(rawEmail)
	if err != nil {
		return nil, err
	}
	phoneNumber, err := value_object.NewPhoneNumber(rawPhoneNumber)
	if err != nil {
		return nil, err
	}
	return &Merchant{
		Id:          id,
		Name:        name,
		Email:       *email,
		PhoneNumber: *phoneNumber,
		Domain:      domain,
		TenantId:    tenantId,
		Status:      Status(status),
	}, nil
}

func (m *Merchant) IsActive() bool {
	return m.Status == ACTIVE
}

func (m *Merchant) Inactive() {
	m.Status = INACTIVE
}
