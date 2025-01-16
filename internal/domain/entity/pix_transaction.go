package entity

import (
	"bom-pedido-api/internal/domain/value_object"
)

const (
	PixTransactionCreated   PixTransactionStatus = "CREATED"
	PixTransactionPaid      PixTransactionStatus = "PAID"
	PixTransactionCancelled PixTransactionStatus = "CANCELLED"
	PixTransactionRefunded  PixTransactionStatus = "REFUNDED"
)

type (
	PixTransactionStatus string

	PixTransaction struct {
		Id             string               `bson:"id"`
		PaymentId      string               `bson:"paymentId"`
		OrderId        string               `bson:"orderId"`
		Status         PixTransactionStatus `bson:"status"`
		Amount         float64              `bson:"amount"`
		PaymentGateway string               `bson:"paymentGateway"`
		QrCode         string               `bson:"qrCode"`
		QrCodeLink     string               `bson:"qrCodeLink"`
	}
)

func NewPixTransaction(paymentId, orderId, qrCode, paymentGateway, qrCodeLink string, amount float64) *PixTransaction {
	return &PixTransaction{
		QrCode:         qrCode,
		QrCodeLink:     qrCodeLink,
		PaymentGateway: paymentGateway,
		Id:             value_object.NewID(),
		PaymentId:      paymentId,
		OrderId:        orderId,
		Status:         PixTransactionCreated,
		Amount:         amount,
	}
}

func (t *PixTransaction) Pay() {
	t.Status = PixTransactionPaid
}

func (t *PixTransaction) Refund() {
	t.Status = PixTransactionRefunded
}

func (t *PixTransaction) Cancel() {
	t.Status = PixTransactionCancelled
}

func (t *PixTransaction) IsPaid() bool {
	return t.Status == PixTransactionPaid
}

func (t *PixTransaction) IsCreated() bool {
	return t.Status == PixTransactionCreated
}

func (t *PixTransaction) IsRefunded() bool {
	return t.Status == PixTransactionRefunded
}

func (t *PixTransaction) IsCancelled() bool {
	return t.Status == PixTransactionCancelled
}
