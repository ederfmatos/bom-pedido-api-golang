package factory

import (
	"bom-pedido-api/application/gateway"
)

type GatewayFactory struct {
	GoogleGateway       gateway.GoogleGateway
	PixGateway          gateway.PixGateway
	NotificationGateway gateway.NotificationGateway
	EmailGateway        gateway.EmailGateway
}

func NewGatewayFactory(
	googleGateway gateway.GoogleGateway,
	pixGateway gateway.PixGateway,
	notificationGateway gateway.NotificationGateway,
	emailGateway gateway.EmailGateway,
) *GatewayFactory {
	return &GatewayFactory{
		GoogleGateway:       googleGateway,
		PixGateway:          pixGateway,
		NotificationGateway: notificationGateway,
		EmailGateway:        emailGateway,
	}
}
