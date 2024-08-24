package factory

import (
	"bom-pedido-api/application/gateway"
)

type GatewayFactory struct {
	GoogleGateway gateway.GoogleGateway
	PixGateway    gateway.PixGateway
}

func NewGatewayFactory(googleGateway gateway.GoogleGateway, pixGateway gateway.PixGateway) *GatewayFactory {
	return &GatewayFactory{GoogleGateway: googleGateway, PixGateway: pixGateway}
}
