package clone

import (
	"bom-pedido-api/internal/application/usecase/product"
	"bom-pedido-api/internal/application/usecase/shopping_cart"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/enums"
	"bom-pedido-api/internal/infra/factory"
	"bom-pedido-api/internal/infra/test"
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
	t.SkipNow()

	ctx := context.Background()
	container := test.NewContainer()
	defer container.Down()
	applicationFactory := factory.NewApplicationFactory(container.GetEnvironment(), container.RedisClient, container.MongoClient)

	merchant, err := entity.NewMerchant(faker.Name(), faker.Email(), faker.Phonenumber(), faker.DomainName())
	require.NoError(t, err)

	err = applicationFactory.MerchantRepository.Create(ctx, merchant)
	require.NoError(t, err)

	customer, err := entity.NewCustomer(faker.Name(), faker.Email(), merchant.TenantId)
	require.NoError(t, err)
	err = applicationFactory.CustomerRepository.Create(ctx, customer)
	require.NoError(t, err)

	category := entity.NewCategory(faker.Name(), faker.Word(), faker.Word())
	err = applicationFactory.ProductCategoryRepository.Create(ctx, category)
	require.NoError(t, err)

	createProduct := product.NewCreateProduct(applicationFactory)
	createProductOutput, err := createProduct.Execute(ctx, product.CreateProductInput{
		Name:        faker.Name(),
		Description: faker.Word(),
		Price:       10.0,
		TenantId:    merchant.TenantId,
		CategoryId:  category.Id,
	})
	require.NoError(t, err)

	addItemToShoppingCart := shopping_cart.NewAddItemToShoppingCart(applicationFactory)
	err = addItemToShoppingCart.Execute(ctx, shopping_cart.AddItemToShoppingCartInput{
		CustomerId:  customer.Id,
		ProductId:   createProductOutput.Id,
		Quantity:    1,
		Observation: "",
		TenantId:    merchant.TenantId,
	})
	require.NoError(t, err)

	checkoutShoppingCart := shopping_cart.NewCheckoutShoppingCart(applicationFactory)
	checkoutOutput, err := checkoutShoppingCart.Execute(ctx, shopping_cart.CheckoutShoppingCartInput{
		CustomerId:      customer.Id,
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

	savedShoppingCart, err := applicationFactory.ShoppingCartRepository.FindByCustomerId(ctx, customer.Id)
	require.NoError(t, err)
	require.NotNil(t, savedShoppingCart)
	require.Equal(t, len(savedShoppingCart.Items), 1)
	for _, shoppingCartItem := range savedShoppingCart.Items {
		require.Equal(t, shoppingCartItem.ProductId, createProductOutput.Id)
		require.Equal(t, shoppingCartItem.Price, 10.0)
		require.Equal(t, shoppingCartItem.Quantity, 1)
		require.Equal(t, shoppingCartItem.Observation, "")
	}
}
