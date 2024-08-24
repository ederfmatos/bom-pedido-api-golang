package gateway

import (
	"context"
	"time"
)

type (
	PixMerchant struct {
		Id    string
		Name  string
		Email string
	}
	CreateQrCodePixInput struct {
		Amount          float64
		InternalOrderId string
		Description     string
		Merchant        PixMerchant
	}
	CreateQrCodePixOutput struct {
		Id             string
		QrCode         string
		ExpiresAt      time.Time
		PaymentGateway string
		QrCodeLink     string
	}
	PixGateway interface {
		CreateQrCodePix(ctx context.Context, input CreateQrCodePixInput) (*CreateQrCodePixOutput, error)
		Name() string
	}
)
