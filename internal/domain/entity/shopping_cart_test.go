package entity

import (
	"bom-pedido-api/internal/domain/enums"
	"bom-pedido-api/internal/domain/errors"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/pkg/faker"
	"bom-pedido-api/pkg/testify/require"
	"fmt"
	"testing"
	"time"
)

func TestShoppingCart_Checkout(t *testing.T) {
	tests := []struct {
		paymentMethod string
		deliveryMode  string
		paymentMode   string
		cardToken     string
		payback       float64
		Errors        []error
	}{
		{
			paymentMethod: "",
			deliveryMode:  "",
			paymentMode:   "",
			payback:       -1,
			Errors:        []error{enums.InvalidPaymentMethodError, enums.InvalidDeliveryModeError, enums.InvalidPaymentModeError, errors.PaybackShouldBePositiveError},
		},
		{
			paymentMethod: enums.CreditCard,
			deliveryMode:  "",
			paymentMode:   enums.PaymentModeInApp.String(),
			payback:       -2,
			cardToken:     "",
			Errors:        []error{enums.InvalidDeliveryModeError, errors.CardTokenIsRequiredError, errors.PaybackShouldBePositiveError},
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("should return errors %v", test.Errors), func(t *testing.T) {
			shoppingCart := NewShoppingCart(value_object.NewID(), faker.Word())

			newProduct, _ := NewProduct(faker.Name(), faker.Word(), 11.0, faker.Word(), faker.Word())
			err := shoppingCart.AddItem(newProduct, 1, faker.Word())
			require.NoError(t, err)

			order, err := shoppingCart.Checkout(test.paymentMethod, test.deliveryMode, test.paymentMode, test.cardToken, test.payback, make(map[string]*Product), time.Second, "")
			require.Nil(t, order)
			expectedError := errors.NewCompositeWithError(test.Errors...)
			require.Equal(t, err, expectedError)
		})
	}

	t.Run("should return product errors", func(t *testing.T) {
		product1, _ := NewProduct(faker.Name(), faker.Word(), 11.0, faker.Word(), faker.Word())
		product2, _ := NewProduct(faker.Name(), faker.Word(), 12.0, faker.Word(), faker.Word())
		product3, _ := NewProduct(faker.Name(), faker.Word(), 13.0, faker.Word(), faker.Word())

		products := map[string]*Product{product1.Id: product1, product2.Id: product2}

		shoppingCart := NewShoppingCart(value_object.NewID(), faker.Word())
		err := shoppingCart.AddItem(product2, 1, faker.Word())
		require.NoError(t, err)
		err = shoppingCart.AddItem(product3, 1, faker.Word())
		require.NoError(t, err)

		product2.MarkUnAvailable()

		order, err := shoppingCart.Checkout(enums.CreditCard, enums.Delivery, enums.InReceiving, "", 0, products, time.Second, "")
		require.Nil(t, order)

		expectedError := errors.NewCompositeWithError(errors.ProductUnAvailableError, errors.ProductNotFoundError)
		require.Equal(t, err, expectedError)
	})

	t.Run("should checkout a shopping cart", func(t *testing.T) {
		product1, _ := NewProduct(faker.Name(), faker.Word(), 11.0, faker.Word(), faker.Word())
		product2, _ := NewProduct(faker.Name(), faker.Word(), 12.0, faker.Word(), faker.Word())
		product3, _ := NewProduct(faker.Name(), faker.Word(), 13.0, faker.Word(), faker.Word())

		products := map[string]*Product{product1.Id: product1, product2.Id: product2, product3.Id: product3}

		shoppingCart := NewShoppingCart(value_object.NewID(), faker.Word())
		err := shoppingCart.AddItem(product1, 1, faker.Word())
		require.NoError(t, err)
		err = shoppingCart.AddItem(product2, 1, faker.Word())
		require.NoError(t, err)

		order, err := shoppingCart.Checkout(enums.CreditCard, enums.Delivery, enums.InReceiving, "", 0, products, time.Second, "")
		require.NoError(t, err)
		require.NotNil(t, order)
		require.Equal(t, shoppingCart.CustomerId, order.CustomerID)
		require.Equal(t, enums.PaymentMethodCreditCard, order.PaymentMethod)
		require.Equal(t, enums.DeliveryModeDelivery, order.DeliveryMode)
		require.Equal(t, enums.PaymentModeInReceiving, order.PaymentMode)
		require.Equal(t, "", order.CreditCardToken)
		require.Equal(t, float64(0), order.Payback)
	})
}
