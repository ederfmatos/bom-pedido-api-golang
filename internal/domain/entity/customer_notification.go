package entity

type CustomerNotification struct {
	CustomerId string `bson:"_id"`
	Recipient  string `bson:"recipient"` // TODO: Mover para dentro do customer
}

func NewCustomerNotification(customerId, recipient string) *CustomerNotification {
	return &CustomerNotification{
		CustomerId: customerId,
		Recipient:  recipient,
	}
}
