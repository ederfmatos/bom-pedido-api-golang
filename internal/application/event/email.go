package event

import (
	"bom-pedido-api/internal/domain/value_object"
)

var (
	SendEmail = "SEND_EMAIL"
)

func NewSendEmailEvent(to value_object.Email, subject string, data map[string]string) *Event {
	emailData := data
	emailData["to"] = to.Value()
	emailData["subject"] = subject
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: to.Value(),
		Name:          SendEmail,
		Data:          emailData,
	}
}
