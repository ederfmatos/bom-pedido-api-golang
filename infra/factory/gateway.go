package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/env"
	"bom-pedido-api/infra/gateway"
)

func gatewayFactory(environment *env.Environment) *factory.GatewayFactory {
	return factory.NewGatewayFactory(
		gateway.NewDefaultGoogleGateway(environment.GoogleAuthUrl),
	)
}
