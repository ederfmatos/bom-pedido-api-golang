package notification

import (
	"bom-pedido-api/domain/value_object"
	"time"
)

type Notification struct {
	Id            string            `bson:"_id"`
	Title         string            `bson:"title"`
	Body          string            `bson:"body"`
	Recipient     string            `bson:"recipient"`
	Status        string            `bson:"status"`
	CorrelationId string            `bson:"correlationId"`
	CreatedAt     time.Time         `bson:"createdAt"`
	Data          map[string]string `bson:"data"`
}

func New(title, body, recipient, correlationId string) *Notification {
	return &Notification{
		Id:            value_object.NewID(),
		Title:         title,
		Body:          body,
		Recipient:     recipient,
		CorrelationId: correlationId,
		Data:          map[string]string{},
		Status:        "CREATED",
		CreatedAt:     time.Now(),
	}
}

func (f *Notification) Put(key, value string) *Notification {
	f.Data[key] = value
	return f
}

func (f *Notification) Fail() {
	f.Status = "ERROR"
}
