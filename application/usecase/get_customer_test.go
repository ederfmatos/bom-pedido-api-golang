package usecase

import (
	"bom-pedido-api/domain/entity"
	"bom-pedido-api/infra/factory"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCustomerUseCase_Execute(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()

	t.Run("it should return CustomerNotFoundError when the customer does not exist", func(t *testing.T) {
		useCase := NewGetCustomerUseCase(applicationFactory)
		input := GetCustomerInput{Id: faker.UUIDDigit(), Context: context.Background()}

		output, err := useCase.Execute(input)

		assert.Nil(t, output)
		assert.ErrorIs(t, entity.CustomerNotFoundError, err)
	})

	t.Run("should return a customer", func(t *testing.T) {
		customer, _ := entity.NewCustomer(faker.Name(), faker.Email())
		_ = customer.SetPhoneNumber(faker.Phonenumber())
		_ = applicationFactory.CustomerRepository.Create(context.Background(), customer)

		useCase := NewGetCustomerUseCase(applicationFactory)
		input := GetCustomerInput{Id: customer.Id, Context: context.Background()}

		output, err := useCase.Execute(input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, customer.Name, output.Name)
		assert.Equal(t, customer.GetEmail(), output.Email)
		assert.Equal(t, customer.GetPhoneNumber(), output.PhoneNumber)
	})
}
