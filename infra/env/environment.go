package env

import "github.com/spf13/viper"

type Environment struct {
	DatabaseUrl                string  `mapstructure:"DATABASE_URL"`
	DatabaseDriver             string  `mapstructure:"DATABASE_DRIVER"`
	RedisUrl                   string  `mapstructure:"REDIS_URL"`
	JwePrivateKeyPath          string  `mapstructure:"JWE_PRIVATE_KEY_PATH"`
	RabbitMqServer             string  `mapstructure:"RABBITMQ_SERVER"`
	GoogleAuthUrl              string  `mapstructure:"GOOGLE_AUTH_URL"`
	AwsClientId                string  `mapstructure:"AWS_CLIENT_ID"`
	AwsClientSecret            string  `mapstructure:"AWS_CLIENT_SECRET"`
	AwsRegion                  string  `mapstructure:"AWS_REGION"`
	AwsEndpoint                *string `mapstructure:"AWS_ENDPOINT"`
	TransactionOutboxTableName string  `mapstructure:"TRANSACTION_OUTBOX_TABLE_NAME"`
}

func LoadEnvironment(path string) *Environment {
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	environment := Environment{}
	err = viper.Unmarshal(&environment)
	if err != nil {
		panic(err)
	}
	return &environment
}
