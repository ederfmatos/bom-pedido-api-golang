package messaging

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/gateway"
	"context"
)

type EmailHandler struct {
	emailGateway gateway.EmailGateway
}

func NewEmailHandler(emailGateway gateway.EmailGateway) *EmailHandler {
	return &EmailHandler{emailGateway: emailGateway}
}

func (h EmailHandler) HandleSendEmail(ctx context.Context, message *event.MessageEvent) error {
	data := message.Event.Data
	err := h.emailGateway.Send(data["to"], data["subject"], data["template"], data)
	return message.AckIfNoError(ctx, err)
}
