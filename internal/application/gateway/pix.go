package gateway

import (
	"bom-pedido-api/internal/domain/errors"
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
		Credential      string
	}
	CreateQrCodePixOutput struct {
		Id             string
		QrCode         string
		ExpiresAt      time.Time
		PaymentGateway string
		QrCodeLink     string
	}
	GetPaymentOutput struct {
		Id             string
		QrCode         string
		ExpiresAt      time.Time
		PaymentGateway string
		QrCodeLink     string
		Status         PaymentStatus
	}
	RefundPixInput struct {
		PaymentId      string
		MerchantId     string
		Amount         float64
		Credential     string
		PaymentGateway string
	}
	GetPaymentInput struct {
		PaymentId      string
		MerchantId     string
		Credential     string
		PaymentGateway string
	}

	PixGateway interface {
		Name() string
		CreateQrCodePix(ctx context.Context, input CreateQrCodePixInput) (*CreateQrCodePixOutput, error)
		RefundPix(ctx context.Context, input RefundPixInput) error
		GetPaymentById(ctx context.Context, input GetPaymentInput) (*GetPaymentOutput, error)
	}
)
