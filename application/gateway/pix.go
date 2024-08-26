package gateway

import (
	"bom-pedido-api/domain/errors"
	"context"
	"time"
)

const (
	TransactionPending   PaymentStatus = "PENDING"
	TransactionPaid      PaymentStatus = "PAID"
	TransactionRefunded  PaymentStatus = "REFUNDED"
	TransactionCancelled PaymentStatus = "CANCELLED"
)

var MerchantGatewayConfigNotFoundError = errors.New("merchant gateway config not found")

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
	RefundPixInput struct {
		PaymentId  string
		MerchantId string
	}

	PixGateway interface {
		CreateQrCodePix(ctx context.Context, input CreateQrCodePixInput) (*CreateQrCodePixOutput, error)
		RefundPix(ctx context.Context, input RefundPixInput) error
		Name() string
		GetPaymentStatus(ctx context.Context, merchantId, id string) (*PaymentStatus, error)
	}
)
