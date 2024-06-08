package token

import (
	"github.com/stretchr/testify/mock"
)

type CustomerTokenManagerMock struct {
	mock.Mock
}

func NewFakeCustomerTokenManager() *CustomerTokenManagerMock {
	return &CustomerTokenManagerMock{}
}

func (tokenManager *CustomerTokenManagerMock) Encrypt(id string) (string, error) {
	args := tokenManager.Called(id)
	return args.String(0), args.Error(1)
}
