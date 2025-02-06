package factory

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/internal/infra/config"
	"bom-pedido-api/internal/infra/gateway/email"
	"bom-pedido-api/internal/infra/gateway/google"
	"bom-pedido-api/internal/infra/gateway/notification"
	"bom-pedido-api/internal/infra/gateway/pix"
	"bom-pedido-api/internal/infra/http_client"
	"bom-pedido-api/internal/infra/repository"
	"bom-pedido-api/pkg/mongo"
	"context"
	firebase "firebase.google.com/go/v4"
	"fmt"
)

func gatewayFactory(environment *config.Environment, mongoDatabase *mongo.Database) (*factory.GatewayFactory, error) {
	paymentGatewayConfigRepository := repository.NewMerchantPaymentGatewayConfigMongoRepository(mongoDatabase)
	pixEnvironment := environment.PixPaymentGateway
	expirationTimeInMinutes := pixEnvironment.ExpirationTimeInMinutes
	pixGateways := map[string]gateway.PixGateway{
		pix.MercadoPago: pix.NewMercadoPagoPixGateway(pixEnvironment.NotificationUrl, expirationTimeInMinutes),
		pix.Woovi:       pix.NewWooviPixGateway(http_client.NewDefaultHttpClient(pixEnvironment.WooviApiBaseUrl), expirationTimeInMinutes),
	}
	for key, pixGateway := range pixGateways {
		pixGateways[key] = pix.NewTelemetryPixGatewayDecorator(pixGateway)
	}
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("create firebase app: %v", err)
	}
	fcmClient, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("create firebase messaging: %v", err)
	}
	return factory.NewGatewayFactory(
		google.NewDefaultGoogleGateway(http_client.NewDefaultHttpClient(environment.GoogleAuthUrl)),
		pix.NewMerchantPixGatewayMediator(paymentGatewayConfigRepository, pixGateways),
		notification.NewTelemetryNotificationGateway(notification.NewFirebaseNotificationGateway(fcmClient)),
		email.NewResendEmailGateway(email.NewTemplateLoader(), environment.EmailFrom, environment.ResendMailKey),
	), nil
}
