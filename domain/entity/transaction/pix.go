package transaction

const (
	CREATED Status = "CREATED"
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
