package enums

type PaymentMode string

const (
	PaymentModeInApp       PaymentMode = "IN_APP"
	PaymentModeInReceiving PaymentMode = "IN_RECEIVING"
)

func (paymentMode PaymentMode) IsValid() bool {
	switch paymentMode {
	case PaymentModeInApp,
		PaymentModeInReceiving:
		return true
	default:
		return false
	}
}

func (paymentMode PaymentMode) IsInApp() bool {
	return paymentMode == PaymentModeInApp
}

func (paymentMode PaymentMode) AllowPaymentMethod(paymentMethod PaymentMethod) bool {
	if !paymentMode.IsValid() || !paymentMethod.IsValid() {
		return false
	}
	if paymentMode == PaymentModeInReceiving {
		return true
	}
	return paymentMethod == PaymentMethodCreditCard
}
