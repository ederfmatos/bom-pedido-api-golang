package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
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

	if err := loadStruct(envStruct); err != nil {
		return nil, fmt.Errorf("load environment struct: %v", err)
	}

	return &environment, nil
}

func loadStruct(s reflect.Value) error {
	for i := 0; i < s.NumField(); i++ {
		field := s.Type().Field(i)
		name := field.Tag.Get("name")

		if name == "" {
			continue
		}

		if field.Type.Kind() == reflect.Struct {
			if err := loadStruct(s.Field(i)); err != nil {
				return fmt.Errorf("load environment struct %s: %v", name, err)
			}
			continue
		}

		value := os.Getenv(name)
		if value == "" {
			return fmt.Errorf("the environment variable %s was not found", name)
		}

		switch s.Field(i).Kind() {
		case reflect.String:
			s.Field(i).SetString(value)
		case reflect.Int:
			intValue, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("failed to convert %s to int: %v", name, err)
			}
			s.Field(i).SetInt(int64(intValue))
		default:
			return fmt.Errorf("unsupported field type: %s", s.Field(i).Kind())
		}
	}
	return nil
}
