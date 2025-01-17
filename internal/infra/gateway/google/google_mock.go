package google

import (
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/pkg/testify/mock"
	"context"
)

type GatewayMock struct {
	mock.Mock
}

func NewFakeGoogleGateway() *GatewayMock {
	return &GatewayMock{}
}

func (googleGateway *GatewayMock) GetUserByToken(_ context.Context, token string) (*gateway.GoogleUser, error) {
	args := googleGateway.Called(token)
	var user = args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*gateway.GoogleUser), args.Error(1)
}
