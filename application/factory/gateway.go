package factory

import (
	"bom-pedido-api/application/gateway"
)

type GatewayFactory struct {
	GoogleGateway gateway.GoogleGateway
}

func NewGatewayFactory(googleGateway gateway.GoogleGateway) *GatewayFactory {
	return &GatewayFactory{GoogleGateway: googleGateway}
}
