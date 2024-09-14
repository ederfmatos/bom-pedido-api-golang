package config

import (
	"log"
	"os"
	"strconv"
)

type (
	PixPaymentGatewayEnv struct {
		NotificationUrl         string
		ExpirationTimeInMinutes int
		WooviApiBaseUrl         string
	}

	Environment struct {
		DatabaseUrl                   string
		DatabaseDriver                string
		RedisUrl                      string
		JwePrivateKeyPath             string
		RabbitMqServer                string
		GoogleAuthUrl                 string
		MongoUrl                      string
		MongoDatabaseName             string
		MongoOutboxCollectionName     string
		Port                          string
		KafkaBootstrapServer          string
		KafkaClientId                 string
		OpenTelemetryEndpointExporter string
		MessagingStrategy             string
		AdminMagicLinkBaseUrl         string
		EmailFrom                     string
		ResendMailKey                 string
		PixPaymentGateway             PixPaymentGatewayEnv
	}
)

func requiredEnv(name string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	log.Fatalf(`Environment variable %s is required`, name)
	return ""
}

func requiredIntEnv(name string) int {
	value, err := strconv.Atoi(requiredEnv(name))
	if err != nil {
		log.Fatal(err)
	}
	return value
}

func optionalEnv(name, defaultValue string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return defaultValue
}

func LoadEnvironment() *Environment {
	return &Environment{
		DatabaseUrl:                   requiredEnv("DATABASE_URL"),
		DatabaseDriver:                requiredEnv("DATABASE_DRIVER"),
		RedisUrl:                      requiredEnv("REDIS_URL"),
		JwePrivateKeyPath:             requiredEnv("JWE_PRIVATE_KEY_PATH"),
		RabbitMqServer:                requiredEnv("RABBITMQ_SERVER"),
		GoogleAuthUrl:                 requiredEnv("GOOGLE_AUTH_URL"),
		MongoUrl:                      requiredEnv("MONGO_URL"),
		MongoDatabaseName:             requiredEnv("MONGO_DATABASE_NAME"),
		MongoOutboxCollectionName:     requiredEnv("MONGO_OUTBOX_COLLECTION"),
		Port:                          optionalEnv("PORT", "8080"),
		KafkaBootstrapServer:          requiredEnv("KAFKA_BOOTSTRAP_SERVER"),
		KafkaClientId:                 requiredEnv("KAFKA_CLIENT_ID"),
		OpenTelemetryEndpointExporter: requiredEnv("OTEL_ENDPOINT_EXPORTER"),
		MessagingStrategy:             requiredEnv("MESSAGING_STRATEGY"),
		AdminMagicLinkBaseUrl:         requiredEnv("ADMIN_MAGIC_LINK_BASE_URL"),
		EmailFrom:                     requiredEnv("EMAIL_FROM"),
		ResendMailKey:                 requiredEnv("RESEND_MAIL_KEY"),
		PixPaymentGateway: PixPaymentGatewayEnv{
			NotificationUrl:         requiredEnv("PIX_PAYMENT_GATEWAY_NOTIFICATION_URL"),
			ExpirationTimeInMinutes: requiredIntEnv("PIX_PAYMENT_GATEWAY_EXPIRATION_TIME_IN_MINUTES"),
			WooviApiBaseUrl:         requiredEnv("WOOVI_API_BASE_URL"),
		},
	}
}
