package customer

type Notification struct {
	CustomerId string `bson:"_id"`
	Recipient  string `bson:"recipient"`
}

func NewNotification(customerId, recipient string) *Notification {
	return &Notification{
		CustomerId: customerId,
		Recipient:  recipient,
	}
}
