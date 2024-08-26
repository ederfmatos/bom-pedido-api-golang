package gateway

import (
	"context"
	"time"
)

const (
	TransactionPending   PaymentStatus = "PENDING"
	TransactionPaid      PaymentStatus = "PAID"
	TransactionRefunded  PaymentStatus = "REFUNDED"
	TransactionCancelled PaymentStatus = "CANCELLED"
)

type (
	PaymentStatus string

	PixCustomer struct {
		Name  string
		Email string
	}
	CreateQrCodePixInput struct {
		Amount          float64
		InternalOrderId string
		Description     string
		MerchantId      string
		Customer        PixCustomer
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
		GetPaymentStatus(ctx context.Context, merchantId, id string) (*PaymentStatus, error)
	}
)
