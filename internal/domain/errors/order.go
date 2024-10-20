package errors

var (
	OrderNotFoundError               = New("order not found")
	OrderDeliveryModeIsWithdrawError = New("the delivery mode for this order is withdraw, delivery not allowed")
	OrderDeliveryModeIsDeliveryError = New("the delivery mode for this order is delivery, withdraw not allowed")
)
