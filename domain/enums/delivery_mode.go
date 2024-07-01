package enums

type DeliveryMode string

const (
	DeliveryModeDelivery DeliveryMode = "DELIVERY"
	DeliveryModeWithdraw DeliveryMode = "WITHDRAW"
)

func (paymentMethod DeliveryMode) IsValid() bool {
	switch paymentMethod {
	case DeliveryModeDelivery,
		DeliveryModeWithdraw:
		return true
	default:
		return false
	}
}
