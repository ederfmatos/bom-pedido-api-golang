package transaction

import (
	"bom-pedido-api/domain/value_object"
)

const (
	created  Status = "CREATED"
	paid     Status = "PAID"
	refunded Status = "REFUNDED"
)

type (
	Status string

	Transaction struct {
		Id        string
		PaymentId string
		OrderId   string
		Status    Status
		Amount    float64
	}

	PixTransaction struct {
		Transaction
		PaymentGateway string
		QrCode         string
		QrCodeLink     string
	}
)

func NewPixTransaction(paymentId, orderId, qrCode, paymentGateway, qrCodeLink string, amount float64) *PixTransaction {
	return &PixTransaction{
		QrCode:         qrCode,
		QrCodeLink:     qrCodeLink,
		PaymentGateway: paymentGateway,
		Transaction: Transaction{
			Id:        value_object.NewID(),
			PaymentId: paymentId,
			OrderId:   orderId,
			Status:    created,
			Amount:    amount,
		},
	}
}

func (t *PixTransaction) Pay() {
	t.Status = paid
}

func (t *PixTransaction) Refund() {
	t.Status = refunded
}

func (t *PixTransaction) IsPaid() bool {
	return t.Status == paid
}

func (t *PixTransaction) IsRefunded() bool {
	return t.Status == refunded
}
