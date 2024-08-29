package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/config"
	"bom-pedido-api/infra/gateway"
	"bom-pedido-api/infra/gateway/pix"
	"bom-pedido-api/infra/repository"
)

func gatewayFactory(environment *config.Environment, connection repository.SqlConnection) *factory.GatewayFactory {
	paymentGatewayConfigRepository := repository.NewDefaultMerchantPaymentGatewayConfigRepository(connection)
	return factory.NewGatewayFactory(
		gateway.NewDefaultGoogleGateway(environment.GoogleAuthUrl),
		pix.NewLogPixGatewayDecorator(
			pix.NewWooviPixGateway(environment.PixPaymentGateway.ExpirationTimeInMinutes, paymentGatewayConfigRepository),
		),
	)
}
