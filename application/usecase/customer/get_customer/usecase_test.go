package get_customer

import (
	"bom-pedido-api/domain/entity/customer"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/infra/factory"
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

		output, err := useCase.Execute(context.TODO(), input)

		require.Nil(t, output)
		require.ErrorIs(t, errors.CustomerNotFoundError, err)
	})

	t.Run("should return a customer", func(t *testing.T) {
		aCustomer, _ := customer.New(faker.Name(), faker.Email(), faker.Word())
		_ = aCustomer.SetPhoneNumber(faker.Phonenumber())
		_ = applicationFactory.CustomerRepository.Create(context.TODO(), aCustomer)

		useCase := New(applicationFactory)
		input := Input{Id: aCustomer.Id}

		output, err := useCase.Execute(context.TODO(), input)

		require.NoError(t, err)
		require.NotNil(t, output)
		require.Equal(t, aCustomer.Name, output.Name)
		require.Equal(t, aCustomer.GetEmail(), output.Email)
		require.Equal(t, aCustomer.GetPhoneNumber(), output.PhoneNumber)
	})
}
