package customer

import (
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/errors"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/factory"
	"bom-pedido-api/pkg/faker"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetCustomerUseCase_Execute(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()

	t.Run("it should return CustomerNotFoundError when the customer does not exist", func(t *testing.T) {
		useCase := NewGetCustomer(applicationFactory)
		input := GetCustomerInput{Id: value_object.NewID()}

		output, err := useCase.Execute(context.Background(), input)

		require.Nil(t, output)
		require.ErrorIs(t, errors.CustomerNotFoundError, err)
	})

	t.Run("should return a customer", func(t *testing.T) {
		customer, _ := entity.NewCustomer(faker.Name(), faker.Email(), faker.Word())
		_ = customer.SetPhoneNumber(faker.PhoneNumber())
		_ = applicationFactory.CustomerRepository.Create(context.Background(), customer)

		useCase := NewGetCustomer(applicationFactory)
		input := GetCustomerInput{Id: customer.Id}

		output, err := useCase.Execute(context.Background(), input)

		require.NoError(t, err)
		require.NotNil(t, output)
		require.Equal(t, customer.Name, output.Name)
		require.Equal(t, customer.GetEmail(), output.Email)
		require.Equal(t, customer.GetPhoneNumber(), output.PhoneNumber)
	})
}
