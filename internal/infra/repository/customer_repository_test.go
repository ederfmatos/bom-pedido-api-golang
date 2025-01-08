package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/customer"
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

			aCustomer, err := customer.New(faker.Name(), faker.Email(), faker.Word())
			require.NoError(t, err)

			savedCustomer, err := customerRepository.FindByEmail(ctx, aCustomer.GetEmail(), aCustomer.TenantId)
			require.NoError(t, err)
			require.Nil(t, savedCustomer)

			savedCustomer, err = customerRepository.FindById(ctx, aCustomer.Id)
			require.NoError(t, err)
			require.Nil(t, savedCustomer)

			err = customerRepository.Create(ctx, aCustomer)
			require.NoError(t, err)

			savedCustomer, err = customerRepository.FindByEmail(ctx, aCustomer.GetEmail(), aCustomer.TenantId)
			require.NoError(t, err)
			require.NotNil(t, savedCustomer)
			require.Equal(t, aCustomer, savedCustomer)

			savedCustomer, err = customerRepository.FindById(ctx, aCustomer.Id)
			require.NoError(t, err)
			require.NotNil(t, savedCustomer)
			require.Equal(t, aCustomer, savedCustomer)

			phoneNumber := "11999999999"
			err = aCustomer.SetPhoneNumber(phoneNumber)
			require.NoError(t, err)

			err = customerRepository.Update(ctx, aCustomer)
			require.NoError(t, err)

			savedCustomer, err = customerRepository.FindByEmail(ctx, aCustomer.GetEmail(), aCustomer.TenantId)
			require.NoError(t, err)
			require.NotNil(t, savedCustomer)
			require.Equal(t, aCustomer, savedCustomer)
			require.Equal(t, phoneNumber, *savedCustomer.GetPhoneNumber())

			savedCustomer, err = customerRepository.FindById(ctx, aCustomer.Id)
			require.NoError(t, err)
			require.NotNil(t, savedCustomer)
			require.Equal(t, aCustomer, savedCustomer)
			require.Equal(t, phoneNumber, *savedCustomer.GetPhoneNumber())
		})
	}
}
