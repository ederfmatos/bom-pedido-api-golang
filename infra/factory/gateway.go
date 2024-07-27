package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/config"
	"bom-pedido-api/infra/gateway"
)

func gatewayFactory(environment *config.Environment) *factory.GatewayFactory {
	return factory.NewGatewayFactory(
		gateway.NewDefaultGoogleGateway(environment.GoogleAuthUrl),
	)
}
