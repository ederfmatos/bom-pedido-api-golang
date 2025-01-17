package token

import (
	"bom-pedido-api/internal/application/token"
	"bom-pedido-api/pkg/testify/mock"
	"context"
)

type CustomerTokenManagerMock struct {
	mock.Mock
}

func NewFakeCustomerTokenManager() *CustomerTokenManagerMock {
	return &CustomerTokenManagerMock{}
}

func (c *CustomerTokenManagerMock) Encrypt(_ context.Context, data token.Data) (string, error) {
	args := c.Called(data)
	return args.String(0), args.Error(1)
}

func (c *CustomerTokenManagerMock) Decrypt(_ context.Context, rawToken string) (*token.Data, error) {
	args := c.Called(rawToken)
	return args.Get(0).(*token.Data), args.Error(1)
}
