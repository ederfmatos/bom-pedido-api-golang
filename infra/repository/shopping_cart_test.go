package repository

import (
	"bom-pedido-api/domain/entity/product"
	"bom-pedido-api/domain/entity/shopping_cart"
	"bom-pedido-api/domain/value_object"
	"context"
	"fmt"
	"github.com/go-faker/faker/v4"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"log"
	"testing"
)

func Test_ShoppingCartRepository(t *testing.T) {
	ctx := context.Background()

	redisContainer, err := redis.Run(ctx, "docker.io/redis:7", redis.WithLogLevel(redis.LogLevelVerbose))
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	defer func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()
	host, err := redisContainer.Host(ctx)
	assert.NoError(t, err)
	port, err := redisContainer.Ports(ctx)
	assert.NoError(t, err)

	options, err := redis2.ParseURL(fmt.Sprintf("redis://%s:%s/0", host, port["6379/tcp"][0].HostPort))
	if err != nil {
		panic(err)
	}
	redisClient := redis2.NewClient(options)
	defer redisClient.Close()

	shoppingCartRepository := NewShoppingCartRedisRepository(redisClient)
	customerId := value_object.NewID()

	shoppingCart, err := shoppingCartRepository.FindByCustomerId(ctx, customerId)
	assert.NoError(t, err)
	assert.Nil(t, shoppingCart)

	shoppingCart = shopping_cart.New(customerId)
	product, err := product.New(faker.Name(), faker.Word(), 10.0)
	assert.NoError(t, err)

	assert.NoError(t, shoppingCart.AddItem(product, 2, ""))
	assert.NoError(t, shoppingCartRepository.Upsert(ctx, shoppingCart))

	savedShoppingCart, err := shoppingCartRepository.FindByCustomerId(ctx, customerId)
	assert.NoError(t, err)
	assert.Equal(t, shoppingCart, savedShoppingCart)

	err = shoppingCartRepository.DeleteByCustomerId(ctx, customerId)
	assert.NoError(t, err)

	savedShoppingCart, err = shoppingCartRepository.FindByCustomerId(ctx, customerId)
	assert.NoError(t, err)
	assert.Nil(t, savedShoppingCart)
}
