package enums

import (
	"bom-pedido-api/internal/domain/errors"
)

type DeliveryMode string

const (
	DeliveryModeDelivery DeliveryMode = "DELIVERY"
	DeliveryModeWithdraw DeliveryMode = "WITHDRAW"
	Delivery                          = string(DeliveryModeDelivery)
	Withdraw                          = string(DeliveryModeWithdraw)
)

var InvalidDeliveryModeError = errors.New("Invalid delivery method")

func ParseDeliveryMode(value string) (DeliveryMode, error) {
	switch value {
	case Delivery:
		return DeliveryModeDelivery, nil
	case Withdraw:
		return DeliveryModeWithdraw, nil
	default:
		return "", InvalidDeliveryModeError
	}
}

func (mode DeliveryMode) String() string {
	return string(mode)
}

func (mode DeliveryMode) IsDelivery() bool {
	return mode == DeliveryModeDelivery
}

func (mode DeliveryMode) IsWithdraw() bool {
	return mode == DeliveryModeWithdraw
}
