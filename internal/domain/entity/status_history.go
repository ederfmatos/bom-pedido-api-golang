package entity

import "time"

type OrderStatusHistory struct {
	Time      time.Time `bson:"time"`
	Status    string    `bson:"status"`
	ChangedBy string    `bson:"changedBy"`
	Data      string    `bson:"data"`
	OrderId   string    `bson:"orderId"`
}

func NewOrderStatusHistory(at time.Time, status, changedBy, data, orderId string) *OrderStatusHistory {
	return &OrderStatusHistory{
		Time:      at.UTC().Truncate(time.Millisecond),
		Status:    status,
		ChangedBy: changedBy,
		Data:      data,
		OrderId:   orderId,
	}
}
