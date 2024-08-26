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
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)

	err = applicationFactory.MerchantRepository.Create(ctx, aMerchant)
	assert.NoError(t, err)

	aCustomer, err := customer.New(faker.Name(), faker.Email(), aMerchant.TenantId)
	assert.NoError(t, err)
	err = applicationFactory.CustomerRepository.Create(ctx, aCustomer)
	assert.NoError(t, err)

	createProduct := create_product.New(applicationFactory)
	createProductOutput, err := createProduct.Execute(ctx, create_product.Input{
		Name:        faker.Name(),
		Description: faker.Word(),
		Price:       10.0,
		TenantId:    aMerchant.TenantId,
	})
	assert.NoError(t, err)

	addItemToShoppingCart := add_item_to_shopping_cart.New(applicationFactory)
	err = addItemToShoppingCart.Execute(ctx, add_item_to_shopping_cart.Input{
		CustomerId:  aCustomer.Id,
		ProductId:   createProductOutput.Id,
		Quantity:    1,
		Observation: "",
		TenantId:    aMerchant.TenantId,
	})
	assert.NoError(t, err)

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
	assert.NoError(t, err)

	orderId := checkoutOutput.Id
	url := fmt.Sprintf("/orders/%s/clone", checkoutOutput.Id)
	request := httptest.NewRequest(http.MethodPost, url, nil)
	response := httptest.NewRecorder()

	c := echo.New().NewContext(request, response)
	c.SetPath("/orders/:id/clone")
	c.SetParamNames("id")
	c.SetParamValues(orderId)
	err = Handle(applicationFactory)(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, response.Code)

	savedShoppingCart, err := applicationFactory.ShoppingCartRepository.FindByCustomerId(ctx, aCustomer.Id)
	assert.NoError(t, err)
	assert.NotNil(t, savedShoppingCart)
	assert.Equal(t, len(savedShoppingCart.Items), 1)
	shoppingCartItem := savedShoppingCart.Items[0]
	assert.Equal(t, shoppingCartItem.ProductId, createProductOutput.Id)
	assert.Equal(t, shoppingCartItem.Price, 10.0)
	assert.Equal(t, shoppingCartItem.Quantity, 1)
	assert.Equal(t, shoppingCartItem.Observation, "")
}
