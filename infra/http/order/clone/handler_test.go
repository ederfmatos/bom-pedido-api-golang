package clone

import (
	"bom-pedido-api/application/usecase/product/create_product"
	"bom-pedido-api/application/usecase/shopping_cart/add_item_to_shopping_cart"
	"bom-pedido-api/application/usecase/shopping_cart/checkout"
	"bom-pedido-api/domain/entity/customer"
	"bom-pedido-api/domain/entity/merchant"
	"bom-pedido-api/domain/enums"
	"bom-pedido-api/infra/factory"
	"bom-pedido-api/infra/test"
	"context"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandle(t *testing.T) {
	ctx := context.Background()
	container := test.NewContainer()
	defer container.Down()
	applicationFactory := factory.NewContainerApplicationFactory(container)

	aMerchant, err := merchant.New(faker.Name(), faker.Email(), faker.Phonenumber(), faker.DomainName())
	require.NoError(t, err)

	err = applicationFactory.MerchantRepository.Create(ctx, aMerchant)
	require.NoError(t, err)

	aCustomer, err := customer.New(faker.Name(), faker.Email(), aMerchant.TenantId)
	require.NoError(t, err)
	err = applicationFactory.CustomerRepository.Create(ctx, aCustomer)
	require.NoError(t, err)

	createProduct := create_product.New(applicationFactory)
	createProductOutput, err := createProduct.Execute(ctx, create_product.Input{
		Name:        faker.Name(),
		Description: faker.Word(),
		Price:       10.0,
		TenantId:    aMerchant.TenantId,
	})
	require.NoError(t, err)

	addItemToShoppingCart := add_item_to_shopping_cart.New(applicationFactory)
	err = addItemToShoppingCart.Execute(ctx, add_item_to_shopping_cart.Input{
		CustomerId:  aCustomer.Id,
		ProductId:   createProductOutput.Id,
		Quantity:    1,
		Observation: "",
		TenantId:    aMerchant.TenantId,
	})
	require.NoError(t, err)

	checkoutShoppingCart := checkout.New(applicationFactory)
	checkoutOutput, err := checkoutShoppingCart.Execute(ctx, checkout.Input{
		CustomerId:      aCustomer.Id,
		PaymentMethod:   enums.Money,
		DeliveryMode:    enums.Withdraw,
		PaymentMode:     enums.InReceiving,
		AddressId:       "",
		Payback:         100,
		CreditCardToken: "",
	})
	require.NoError(t, err)

	orderId := checkoutOutput.Id
	url := fmt.Sprintf("/orders/%s/clone", checkoutOutput.Id)
	request := httptest.NewRequest(http.MethodPost, url, nil)
	response := httptest.NewRecorder()

	c := echo.New().NewContext(request, response)
	c.SetPath("/orders/:id/clone")
	c.SetParamNames("id")
	c.SetParamValues(orderId)
	err = Handle(applicationFactory)(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, response.Code)

	savedShoppingCart, err := applicationFactory.ShoppingCartRepository.FindByCustomerId(ctx, aCustomer.Id)
	require.NoError(t, err)
	require.NotNil(t, savedShoppingCart)
	require.Equal(t, len(savedShoppingCart.Items), 1)
	shoppingCartItem := savedShoppingCart.Items[0]
	require.Equal(t, shoppingCartItem.ProductId, createProductOutput.Id)
	require.Equal(t, shoppingCartItem.Price, 10.0)
	require.Equal(t, shoppingCartItem.Quantity, 1)
	require.Equal(t, shoppingCartItem.Observation, "")
}
