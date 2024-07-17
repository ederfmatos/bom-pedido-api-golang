package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/product"
	"bom-pedido-api/domain/entity/shopping_cart"
	"bom-pedido-api/domain/value_object"
	"context"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func Test_ShoppingCartMongoRepository(t *testing.T) {
	mongoDatabase, closeDatabase := MongoConnection(t)
	defer closeDatabase()
	shoppingCartRepository := NewShoppingCartMongoRepository(mongoDatabase)
	runTests(t, shoppingCartRepository)
}

func Test_ShoppingCartMemoryRepository(t *testing.T) {
	shoppingCartRepository := NewShoppingCartMemoryRepository()
	runTests(t, shoppingCartRepository)
}

func runTests(t *testing.T, shoppingCartRepository repository.ShoppingCartRepository) {
	customerId := value_object.NewID()
	ctx := context.Background()

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

func MongoConnection(t *testing.T) (*mongo.Database, func()) {
	ctx := context.Background()
	mongodbContainer, err := mongodb.Run(ctx, "mongo:6")

	endpoint, err := mongodbContainer.Endpoint(context.Background(), "")
	if err != nil {
		assert.NoError(t, err)
	}

	uri := fmt.Sprintf("mongodb://%s", endpoint)

	clientOptions := options.Client().ApplyURI(uri)
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		assert.NoError(t, err)
	}

	database := mongoClient.Database("test")
	return database, func() {
		go mongodbContainer.Terminate(ctx)
		go mongoClient.Disconnect(ctx)
	}
}
