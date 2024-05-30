package token

import (
	"bom-pedido-api/application/token"
	"github.com/stretchr/testify/mock"
)

type CustomerTokenManagerMock struct {
	mock.Mock
}

func NewFakeCustomerTokenManager() token.CustomerTokenManager {
	return &CustomerTokenManagerMock{}
}

func (tokenManager *CustomerTokenManagerMock) Encrypt(id string) (string, error) {
	args := tokenManager.Called(id)
	return args.String(0), args.Error(1)
}
