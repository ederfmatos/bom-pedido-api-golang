package config

import (
	"os"
)

type Environment struct {
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
}

func requiredEnv(name string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	panic(`Environment variable ` + name + ` is required`)
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
	}
}
