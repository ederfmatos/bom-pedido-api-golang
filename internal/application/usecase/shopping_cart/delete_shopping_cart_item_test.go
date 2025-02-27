package shopping_cart

import (
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/factory"
	"bom-pedido-api/pkg/faker"
	"bom-pedido-api/pkg/testify/require"
	"context"
	"testing"
)

func Test_DeleteShoppingCartItem(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()
	useCase := NewDeleteShoppingCartItem(applicationFactory)
	customerId := value_object.NewID()
	shoppingCart := entity.NewShoppingCart(customerId, faker.Word())

	product, err := entity.NewProduct(faker.Name(), faker.Word(), 10.0, faker.Word(), faker.Word())
	require.NoError(t, err)

	err = shoppingCart.AddItem(product, 1, "")
	require.NoError(t, err)
	var itemId string
	for id := range shoppingCart.Items {
		itemId = id
	}

	input := DeleteShoppingCartItemInput{CustomerId: customerId, ItemId: itemId}
	ctx := context.Background()

	err = applicationFactory.ShoppingCartRepository.Upsert(ctx, shoppingCart)
	require.NoError(t, err)

	err = useCase.Execute(ctx, input)
	require.NoError(t, err)
	savedShoppingCart, err := applicationFactory.ShoppingCartRepository.FindByCustomerId(ctx, customerId)
	require.NoError(t, err)
	require.Equal(t, 0, len(savedShoppingCart.Items))
}
