package shopping_cart

import (
	"bom-pedido-api/domain/entity/product"
	"bom-pedido-api/domain/enums"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/domain/value_object"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShoppingCart_Checkout(t *testing.T) {
	tests := []struct {
		paymentMethod string
		deliveryMode  string
		paymentMode   string
		cardToken     string
		change        float64
		Errors        []error
	}{
		{
			paymentMethod: "",
			deliveryMode:  "",
			paymentMode:   "",
			change:        -1,
			Errors:        []error{enums.InvalidPaymentMethodError, enums.InvalidDeliveryModeError, enums.InvalidPaymentModeError, errors.ChangeShouldBePositiveError},
		},
		{
			paymentMethod: enums.CreditCard,
			deliveryMode:  "",
			paymentMode:   enums.PaymentModeInApp.String(),
			change:        -2,
			cardToken:     "",
			Errors:        []error{enums.InvalidDeliveryModeError, errors.CardTokenIsRequiredError, errors.ChangeShouldBePositiveError},
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("should return errors %v", test.Errors), func(t *testing.T) {
			shoppingCart := New(value_object.NewID())

			newProduct, _ := product.New(faker.Name(), faker.Word(), 11.0)
			err := shoppingCart.AddItem(newProduct, 1, faker.Word())
			assert.NoError(t, err)

			order, err := shoppingCart.Checkout(test.paymentMethod, test.deliveryMode, test.paymentMode, test.cardToken, test.change, make(map[string]*product.Product), time.Second)
			assert.Nil(t, order)
			expectedError := errors.NewCompositeWithError(test.Errors...)
			assert.Equal(t, err, expectedError)
		})
	}

	t.Run("should return product errors", func(t *testing.T) {
		product1, _ := product.New(faker.Name(), faker.Word(), 11.0)
		product2, _ := product.New(faker.Name(), faker.Word(), 12.0)
		product3, _ := product.New(faker.Name(), faker.Word(), 13.0)

		products := map[string]*product.Product{product1.Id: product1, product2.Id: product2}

		shoppingCart := New(value_object.NewID())
		err := shoppingCart.AddItem(product2, 1, faker.Word())
		assert.NoError(t, err)
		err = shoppingCart.AddItem(product3, 1, faker.Word())
		assert.NoError(t, err)

		product2.MarkUnAvailable()

		order, err := shoppingCart.Checkout(enums.CreditCard, enums.Delivery, enums.InReceiving, "", 0, products, time.Second)
		assert.Nil(t, order)

		expectedError := errors.NewCompositeWithError(errors.ProductUnAvailableError, errors.ProductNotFoundError)
		assert.Equal(t, err, expectedError)
	})

	t.Run("should checkout a shopping cart", func(t *testing.T) {
		product1, _ := product.New(faker.Name(), faker.Word(), 11.0)
		product2, _ := product.New(faker.Name(), faker.Word(), 12.0)
		product3, _ := product.New(faker.Name(), faker.Word(), 13.0)

		products := map[string]*product.Product{product1.Id: product1, product2.Id: product2, product3.Id: product3}

		shoppingCart := New(value_object.NewID())
		err := shoppingCart.AddItem(product1, 1, faker.Word())
		assert.NoError(t, err)
		err = shoppingCart.AddItem(product2, 1, faker.Word())
		assert.NoError(t, err)

		order, err := shoppingCart.Checkout(enums.CreditCard, enums.Delivery, enums.InReceiving, "", 0, products, time.Second)
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, shoppingCart.CustomerId, order.CustomerID)
		assert.Equal(t, enums.PaymentMethodCreditCard, order.PaymentMethod)
		assert.Equal(t, enums.DeliveryModeDelivery, order.DeliveryMode)
		assert.Equal(t, enums.PaymentModeInReceiving, order.PaymentMode)
		assert.Equal(t, "", order.CreditCardToken)
		assert.Equal(t, float64(0), order.Change)
	})
}
