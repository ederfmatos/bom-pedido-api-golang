package token

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type CustomerTokenManagerMock struct {
	mock.Mock
}

func NewFakeCustomerTokenManager() *CustomerTokenManagerMock {
	return &CustomerTokenManagerMock{}
}

func (tokenManager *CustomerTokenManagerMock) Encrypt(_ context.Context, id string) (string, error) {
	args := tokenManager.Called(id)
	return args.String(0), args.Error(1)
}

func (tokenManager *CustomerTokenManagerMock) Decrypt(_ context.Context, token string) (string, error) {
	args := tokenManager.Called(token)
	return args.String(0), args.Error(1)
}
