package env

import (
	"github.com/spf13/viper"
	"os"
)

type Environment struct {
	DatabaseUrl               string `mapstructure:"DATABASE_URL"`
	DatabaseDriver            string `mapstructure:"DATABASE_DRIVER"`
	RedisUrl                  string `mapstructure:"REDIS_URL"`
	JwePrivateKeyPath         string `mapstructure:"JWE_PRIVATE_KEY_PATH"`
	RabbitMqServer            string `mapstructure:"RABBITMQ_SERVER"`
	GoogleAuthUrl             string `mapstructure:"GOOGLE_AUTH_URL"`
	MongoUrl                  string `mapstructure:"MONGO_URL"`
	MongoDatabaseName         string `mapstructure:"MONGO_DATABASE_NAME"`
	MongoOutboxCollectionName string `mapstructure:"MONGO_OUTBOX_COLLECTION_NAME"`
	Port                      string `mapstructure:"PORT"`
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
	if environment.Port == "" {
		environment.Port = os.Getenv("PORT")
	}
	if environment.Port == "" {
		environment.Port = "8080"
	}
	return &environment
}
