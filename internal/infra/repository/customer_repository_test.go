package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/infra/test"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_CustomerRepository(t *testing.T) {
	container := test.NewContainer()
	repositories := map[string]repository.CustomerRepository{
		"CustomerMemoryRepository": NewCustomerMemoryRepository(),
		"CustomerMongoRepository":  NewCustomerMongoRepository(container.MongoDatabase()),
	}

	for name, customerRepository := range repositories {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			customer, err := entity.NewCustomer(faker.Name(), faker.Email(), faker.Word())
			require.NoError(t, err)

			savedCustomer, err := customerRepository.FindByEmail(ctx, customer.GetEmail(), customer.TenantId)
			require.NoError(t, err)
			require.Nil(t, savedCustomer)

			savedCustomer, err = customerRepository.FindById(ctx, customer.Id)
			require.NoError(t, err)
			require.Nil(t, savedCustomer)

			err = customerRepository.Create(ctx, customer)
			require.NoError(t, err)

			savedCustomer, err = customerRepository.FindByEmail(ctx, customer.GetEmail(), customer.TenantId)
			require.NoError(t, err)
			require.NotNil(t, savedCustomer)
			require.Equal(t, customer, savedCustomer)

			savedCustomer, err = customerRepository.FindById(ctx, customer.Id)
			require.NoError(t, err)
			require.NotNil(t, savedCustomer)
			require.Equal(t, customer, savedCustomer)

			phoneNumber := "11999999999"
			err = customer.SetPhoneNumber(phoneNumber)
			require.NoError(t, err)

			err = customerRepository.Update(ctx, customer)
			require.NoError(t, err)

			savedCustomer, err = customerRepository.FindByEmail(ctx, customer.GetEmail(), customer.TenantId)
			require.NoError(t, err)
			require.NotNil(t, savedCustomer)
			require.Equal(t, customer, savedCustomer)
			require.Equal(t, phoneNumber, *savedCustomer.GetPhoneNumber())

			savedCustomer, err = customerRepository.FindById(ctx, customer.Id)
			require.NoError(t, err)
			require.NotNil(t, savedCustomer)
			require.Equal(t, customer, savedCustomer)
			require.Equal(t, phoneNumber, *savedCustomer.GetPhoneNumber())
		})
	}
}
