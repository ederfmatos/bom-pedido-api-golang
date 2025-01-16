package get_customer

import (
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/errors"
	"bom-pedido-api/internal/infra/factory"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetCustomerUseCase_Execute(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()

	t.Run("it should return CustomerNotFoundError when the customer does not exist", func(t *testing.T) {
		useCase := New(applicationFactory)
		input := Input{Id: faker.UUIDDigit()}

		output, err := useCase.Execute(context.Background(), input)

		require.Nil(t, output)
		require.ErrorIs(t, errors.CustomerNotFoundError, err)
	})

	t.Run("should return a customer", func(t *testing.T) {
		customer, _ := entity.NewCustomer(faker.Name(), faker.Email(), faker.Word())
		_ = customer.SetPhoneNumber(faker.Phonenumber())
		_ = applicationFactory.CustomerRepository.Create(context.Background(), customer)

		useCase := New(applicationFactory)
		input := Input{Id: customer.Id}

		output, err := useCase.Execute(context.Background(), input)

		require.NoError(t, err)
		require.NotNil(t, output)
		require.Equal(t, customer.Name, output.Name)
		require.Equal(t, customer.GetEmail(), output.Email)
		require.Equal(t, customer.GetPhoneNumber(), output.PhoneNumber)
	})
}
