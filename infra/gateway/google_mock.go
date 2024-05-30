package gateway

import (
	"bom-pedido-api/application/gateway"
	"github.com/stretchr/testify/mock"
)

type GoogleGatewayMock struct {
	mock.Mock
}

func NewFakeGoogleGateway() *GoogleGatewayMock {
	return &GoogleGatewayMock{}
}

func (googleGateway *GoogleGatewayMock) GetUserByToken(token string) (*gateway.GoogleUser, error) {
	args := googleGateway.Called(token)
	var user = args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*gateway.GoogleUser), args.Error(1)
}
