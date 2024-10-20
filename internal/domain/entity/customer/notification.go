package customer

type Notification struct {
	CustomerId string `bson:"_id"`
	Recipient  string `bson:"recipient"`
}
