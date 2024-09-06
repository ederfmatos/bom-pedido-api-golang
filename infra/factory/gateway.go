package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/gateway"
	"bom-pedido-api/infra/config"
	"bom-pedido-api/infra/gateway/google"
	"bom-pedido-api/infra/gateway/notification"
	"bom-pedido-api/infra/gateway/pix"
	"bom-pedido-api/infra/http_client"
	"bom-pedido-api/infra/repository"
	"context"
	firebase "firebase.google.com/go"
)

func gatewayFactory(environment *config.Environment, connection repository.SqlConnection) *factory.GatewayFactory {
	paymentGatewayConfigRepository := repository.NewDefaultMerchantPaymentGatewayConfigRepository(connection)
	pixEnvironment := environment.PixPaymentGateway
	expirationTimeInMinutes := pixEnvironment.ExpirationTimeInMinutes
	pixGateways := map[string]gateway.PixGateway{
		pix.MercadoPago: pix.NewMercadoPagoPixGateway(pixEnvironment.NotificationUrl, expirationTimeInMinutes),
		pix.Woovi:       pix.NewWooviPixGateway(http_client.NewDefaultHttpClient(pixEnvironment.WooviApiBaseUrl), expirationTimeInMinutes),
	}
	for key, pixGateway := range pixGateways {
		pixGateways[key] = pix.NewLogPixGatewayDecorator(pixGateway)
	}
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		panic(err)
	}
	fcmClient, err := app.Messaging(ctx)
	if err != nil {
		panic(err)
	}
	return factory.NewGatewayFactory(
		google.NewDefaultGoogleGateway(http_client.NewDefaultHttpClient(environment.GoogleAuthUrl)),
		pix.NewMerchantPixGatewayMediator(paymentGatewayConfigRepository, pixGateways),
		notification.NewFirebaseNotificationGateway(fcmClient),
	)
}
