package enums

import "bom-pedido-api/domain/errors"

type PaymentMode string

const (
	PaymentModeInApp       PaymentMode = "IN_APP"
	PaymentModeInReceiving PaymentMode = "IN_RECEIVING"
	InApp                              = string(PaymentModeInApp)
	InReceiving                        = string(PaymentModeInReceiving)
)

var AllPaymentModes = []PaymentMode{
	PaymentModeInApp,
	PaymentModeInReceiving,
}

var InvalidPaymentModeError = errors.New("Invalid payment mode")

func ParsePaymentMode(value string) (PaymentMode, error) {
	switch value {
	case InApp:
		return PaymentModeInApp, nil
	case InReceiving:
		return PaymentModeInReceiving, nil
	default:
		return "", InvalidPaymentModeError
	}
}

func (paymentMode PaymentMode) IsInApp() bool {
	return paymentMode == PaymentModeInApp
}

func (paymentMode PaymentMode) String() string {
	return string(paymentMode)
}
