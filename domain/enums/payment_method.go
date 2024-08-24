package enums

import "bom-pedido-api/domain/errors"

type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "CREDIT_CARD"
	PaymentMethodDebitCard  PaymentMethod = "DEBIT_CARD"
	PaymentMethodPix        PaymentMethod = "PIX"
	PaymentMethodMoney      PaymentMethod = "MONEY"
	CreditCard                            = string(PaymentMethodCreditCard)
	DebitCard                             = string(PaymentMethodDebitCard)
	Pix                                   = string(PaymentMethodPix)
	Money                                 = string(PaymentMethodMoney)
)

var AllPaymentMethods = []PaymentMethod{
	PaymentMethodCreditCard,
	PaymentMethodDebitCard,
	PaymentMethodPix,
	PaymentMethodMoney,
}

var InvalidPaymentMethodError = errors.New("Invalid payment method")

func ParsePaymentMethod(value string) (PaymentMethod, error) {
	switch value {
	case CreditCard:
		return PaymentMethodCreditCard, nil
	case DebitCard:
		return PaymentMethodDebitCard, nil
	case Pix:
		return PaymentMethodPix, nil
	case Money:
		return PaymentMethodMoney, nil
	default:
		return "", InvalidPaymentMethodError
	}
}

func (paymentMethod PaymentMethod) IsCreditCard() bool {
	return paymentMethod == PaymentMethodCreditCard
}

func (paymentMethod PaymentMethod) String() string {
	return string(paymentMethod)
}

func (paymentMethod PaymentMethod) IsPix() bool {
	return paymentMethod == PaymentMethodPix
}
