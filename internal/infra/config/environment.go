package config

import (
	"bom-pedido-api/pkg/config"
	"fmt"
	"reflect"
)

type (
	PixPaymentGatewayEnv struct {
		NotificationUrl         string `name:"PIX_PAYMENT_GATEWAY_NOTIFICATION_URL"`
		ExpirationTimeInMinutes int    `name:"PIX_PAYMENT_GATEWAY_EXPIRATION_TIME_IN_MINUTES"`
		WooviApiBaseUrl         string `name:"WOOVI_API_BASE_URL"`
	}

	Environment struct {
		RedisUrl                      string               `name:"REDIS_URL"`
		JwePrivateKeyPath             string               `name:"JWE_PRIVATE_KEY_PATH"`
		RabbitMqServer                string               `name:"RABBITMQ_SERVER"`
		GoogleAuthUrl                 string               `name:"GOOGLE_AUTH_URL"`
		MongoUrl                      string               `name:"MONGO_URL"`
		MongoDatabaseName             string               `name:"MONGO_DATABASE_NAME"`
		MongoOutboxCollectionName     string               `name:"MONGO_OUTBOX_COLLECTION"`
		Port                          string               `name:"PORT"`
		OpenTelemetryEndpointExporter string               `name:"OTEL_ENDPOINT_EXPORTER"`
		MessagingStrategy             string               `name:"MESSAGING_STRATEGY"`
		AdminMagicLinkBaseUrl         string               `name:"ADMIN_MAGIC_LINK_BASE_URL"`
		EmailFrom                     string               `name:"EMAIL_FROM"`
		ResendMailKey                 string               `name:"RESEND_MAIL_KEY"`
		PixPaymentGateway             PixPaymentGatewayEnv `name:"PIX_PAYMENT_GATEWAY"`
	}
)

func LoadEnvironment() (*Environment, error) {
	var environment Environment
	envStruct := reflect.ValueOf(&environment).Elem()

	if err := config.Load(envStruct); err != nil {
		return nil, fmt.Errorf("load environment struct: %v", err)
	}

	return &environment, nil
}
