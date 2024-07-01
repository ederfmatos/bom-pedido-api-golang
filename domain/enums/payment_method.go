package enums

type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "CREDIT_CARD"
	PaymentMethodDebitCard  PaymentMethod = "DEBIT_CARD"
	PaymentMethodPix        PaymentMethod = "PIX"
	PaymentMethodMoney      PaymentMethod = "MONEY"
)

func (paymentMethod PaymentMethod) IsValid() bool {
	switch paymentMethod {
	case PaymentMethodCreditCard,
		PaymentMethodDebitCard,
		PaymentMethodPix,
		PaymentMethodMoney:
		return true
	default:
		return false
	}
}

func (paymentMethod PaymentMethod) IsCreditCard() bool {
	return paymentMethod == PaymentMethodCreditCard
}
