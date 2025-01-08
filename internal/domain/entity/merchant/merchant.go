package merchant

import (
	"bom-pedido-api/internal/domain/value_object"
)

const (
	ACTIVE   Status = "ACTIVE"
	INACTIVE Status = "INACTIVE"
)

type (
	Status string

	Merchant struct {
		Id          string                   `bson:"id"`
		Name        string                   `bson:"name"`
		Email       value_object.Email       `bson:"email"`
		PhoneNumber value_object.PhoneNumber `bson:"phoneNumber"`
		Domain      string                   `bson:"domain"`
		TenantId    string                   `bson:"tenantId"`
		Status      Status                   `bson:"status"`
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

func (m *Merchant) IsActive() bool {
	return m.Status == ACTIVE
}

func (m *Merchant) Inactive() {
	m.Status = INACTIVE
}
