package email

import (
	"bom-pedido-api/internal/application/gateway"
)

type FakeEmailGateway struct {
}

func NewFakeEmailGateway() gateway.EmailGateway {
	return &FakeEmailGateway{}
}

func (r *FakeEmailGateway) Send(string, string, string, map[string]string) error {
	return nil
}
