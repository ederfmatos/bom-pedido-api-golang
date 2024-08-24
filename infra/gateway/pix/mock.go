package pix

import (
	"bom-pedido-api/application/gateway"
	"context"
	"github.com/stretchr/testify/mock"
)

type MockPixGateway struct {
	mock.Mock
}

func (f *MockPixGateway) Name() string {
	return "MOCK"
}

func NewFakePixGateway() *MockPixGateway {
	return &MockPixGateway{}
}

func (f *MockPixGateway) CreateQrCodePix(ctx context.Context, input gateway.CreateQrCodePixInput) (*gateway.CreateQrCodePixOutput, error) {
	args := f.Called(ctx, input)
	var output = args.Get(0)
	if output == nil {
		return nil, args.Error(1)
	}
	return output.(*gateway.CreateQrCodePixOutput), args.Error(1)
}
