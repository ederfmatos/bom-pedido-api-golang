package transaction

import (
	"bom-pedido-api/internal/domain/value_object"
)

const (
	created   Status = "CREATED"
	paid      Status = "PAID"
	cancelled Status = "CANCELLED"
	refunded  Status = "REFUNDED"
)

type (
	Status string

	PixTransaction struct {
		Id             string  `bson:"id"`
		PaymentId      string  `bson:"paymentId"`
		OrderId        string  `bson:"orderId"`
		Status         Status  `bson:"status"`
		Amount         float64 `bson:"amount"`
		PaymentGateway string  `bson:"paymentGateway"`
		QrCode         string  `bson:"qrCode"`
		QrCodeLink     string  `bson:"qrCodeLink"`
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
		Status:         created,
		Amount:         amount,
	}
}

func (t *PixTransaction) Pay() {
	t.Status = paid
}

func (t *PixTransaction) Refund() {
	t.Status = refunded
}

func (t *PixTransaction) Cancel() {
	t.Status = cancelled
}

func (t *PixTransaction) IsPaid() bool {
	return t.Status == paid
}

func (t *PixTransaction) IsCreated() bool {
	return t.Status == created
}

func (t *PixTransaction) IsRefunded() bool {
	return t.Status == refunded
}

func (t *PixTransaction) IsCancelled() bool {
	return t.Status == cancelled
}
