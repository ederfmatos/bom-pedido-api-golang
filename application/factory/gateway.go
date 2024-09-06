package factory

import (
	"bom-pedido-api/application/gateway"
)

type GatewayFactory struct {
	GoogleGateway       gateway.GoogleGateway
	PixGateway          gateway.PixGateway
	NotificationGateway gateway.NotificationGateway
}

func NewGatewayFactory(
	googleGateway gateway.GoogleGateway,
	pixGateway gateway.PixGateway,
	notificationGateway gateway.NotificationGateway,
) *GatewayFactory {
	return &GatewayFactory{
		GoogleGateway:       googleGateway,
		PixGateway:          pixGateway,
		NotificationGateway: notificationGateway,
	}
}
