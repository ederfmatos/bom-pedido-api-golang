package pix

import (
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/pkg/testify/mock"
	"context"
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

func (f *MockPixGateway) GetPaymentById(ctx context.Context, input gateway.GetPaymentInput) (*gateway.GetPaymentOutput, error) {
	args := f.Called(ctx, input)
	var output = args.Get(0)
	if output == nil {
		return nil, args.Error(1)
	}
	return output.(*gateway.GetPaymentOutput), args.Error(1)
}

func (f *MockPixGateway) RefundPix(ctx context.Context, input gateway.RefundPixInput) error {
	args := f.Called(ctx, input)
	return args.Error(0)
}
