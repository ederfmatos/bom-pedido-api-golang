package order

import "time"

type StatusHistory struct {
	Time      time.Time
	Status    string
	ChangedBy string
	Data      string
	OrderId   string
}
