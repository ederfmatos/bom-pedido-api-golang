package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/test"
	"bom-pedido-api/pkg/faker"
	"bom-pedido-api/pkg/testify/require"
	"context"
	"testing"
)

func Test_ShoppingCartMongoRepository(t *testing.T) {
	container := test.NewContainer()

	repositories := map[string]repository.ShoppingCartRepository{
		"ShoppingCartMemoryRepository": NewShoppingCartMemoryRepository(),
		"ShoppingCartMongoRepository":  NewShoppingCartMongoRepository(container.MongoDatabase()),
	}

	for name, shoppingCartRepository := range repositories {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			customerId := value_object.NewID()
			ctx := context.Background()

			shoppingCart, err := shoppingCartRepository.FindByCustomerId(ctx, customerId)
			require.NoError(t, err)
			require.Nil(t, shoppingCart)

			shoppingCart = entity.NewShoppingCart(customerId, faker.Word())
			product, err := entity.NewProduct(faker.Name(), faker.Word(), 10.0, faker.Word(), faker.Word())
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
		})
	}
}
