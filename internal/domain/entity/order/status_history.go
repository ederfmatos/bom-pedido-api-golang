package order

import "time"

type StatusHistory struct {
	Time      time.Time `bson:"time"`
	Status    string    `bson:"status"`
	ChangedBy string    `bson:"changedBy"`
	Data      string    `bson:"data"`
	OrderId   string    `bson:"orderId"`
}

func NewStatusHistory(at time.Time, status string, changedBy string, data string, orderId string) *StatusHistory {
	return &StatusHistory{
		Time:      at.UTC().Truncate(time.Millisecond),
		Status:    status,
		ChangedBy: changedBy,
		Data:      data,
		OrderId:   orderId,
	}
}
