package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/product"
	"bom-pedido-api/domain/entity/shopping_cart"
	"bom-pedido-api/domain/value_object"
	"bom-pedido-api/infra/test"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ShoppingCartMongoRepository(t *testing.T) {
	container := test.NewContainer()
	shoppingCartRepository := NewShoppingCartMongoRepository(container.MongoDatabase)
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
	require.NoError(t, err)
	require.Nil(t, shoppingCart)

	shoppingCart = shopping_cart.New(customerId, faker.WORD)
	product, err := product.New(faker.Name(), faker.Word(), 10.0, faker.WORD)
	require.NoError(t, err)

	require.NoError(t, shoppingCart.AddItem(product, 2, ""))
	require.NoError(t, shoppingCartRepository.Upsert(ctx, shoppingCart))

	savedShoppingCart, err := shoppingCartRepository.FindByCustomerId(ctx, customerId)
	require.NoError(t, err)
	require.Equal(t, shoppingCart, savedShoppingCart)

	err = shoppingCartRepository.DeleteByCustomerId(ctx, customerId)
	require.NoError(t, err)

	savedShoppingCart, err = shoppingCartRepository.FindByCustomerId(ctx, customerId)
	require.NoError(t, err)
	require.Nil(t, savedShoppingCart)
}
