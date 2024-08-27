package transaction

import "bom-pedido-api/domain/errors"

const (
	CREATED Status = "CREATED"
	PAID    Status = "PAID"
)

var (
	AlreadyPaidError = errors.New("Transaction already paid")
)

type (
	Status string

	Transaction struct {
		Id      string
		OrderId string
		Status  Status
		Amount  float64
	}

	PixTransaction struct {
		Transaction
		PaymentGateway string
		QrCode         string
		QrCodeLink     string
	}
)

func NewPixTransaction(id, orderId, qrCode, paymentGateway, qrCodeLink string, amount float64) *PixTransaction {
	return &PixTransaction{
		QrCode:         qrCode,
		QrCodeLink:     qrCodeLink,
		PaymentGateway: paymentGateway,
		Transaction: Transaction{
			Id:      id,
			OrderId: orderId,
			Status:  CREATED,
			Amount:  amount,
		},
	}
}

func (t *PixTransaction) Pay() error {
	if t.Status == PAID {
		return AlreadyPaidError
	}
	t.Status = PAID
	return nil
}
