package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/customer"
	"bom-pedido-api/infra/test"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_CustomerSqlRepository(t *testing.T) {
	container := test.NewContainer()
	sqlConnection := NewDefaultSqlConnection(container.Database)
	customerSqlRepository := NewDefaultCustomerRepository(sqlConnection)
	runCustomerTests(t, customerSqlRepository)
}

func Test_CustomerMemoryRepository(t *testing.T) {
	customerSqlRepository := NewCustomerMemoryRepository()
	runCustomerTests(t, customerSqlRepository)
}

func runCustomerTests(t *testing.T, repository repository.CustomerRepository) {
	ctx := context.Background()

	aCustomer, err := customer.New(faker.Name(), faker.Email(), faker.Word())
	require.NoError(t, err)

	savedCustomer, err := repository.FindByEmail(ctx, aCustomer.GetEmail(), aCustomer.TenantId)
	require.NoError(t, err)
	require.Nil(t, savedCustomer)

	savedCustomer, err = repository.FindById(ctx, aCustomer.Id)
	require.NoError(t, err)
	require.Nil(t, savedCustomer)

	err = repository.Create(ctx, aCustomer)
	require.NoError(t, err)

	savedCustomer, err = repository.FindByEmail(ctx, aCustomer.GetEmail(), aCustomer.TenantId)
	require.NoError(t, err)
	require.NotNil(t, savedCustomer)
	require.Equal(t, aCustomer, savedCustomer)

	savedCustomer, err = repository.FindById(ctx, aCustomer.Id)
	require.NoError(t, err)
	require.NotNil(t, savedCustomer)
	require.Equal(t, aCustomer, savedCustomer)

	phoneNumber := "11999999999"
	err = aCustomer.SetPhoneNumber(phoneNumber)
	require.NoError(t, err)

	err = repository.Update(ctx, aCustomer)
	require.NoError(t, err)

	savedCustomer, err = repository.FindByEmail(ctx, aCustomer.GetEmail(), aCustomer.TenantId)
	require.NoError(t, err)
	require.NotNil(t, savedCustomer)
	require.Equal(t, aCustomer, savedCustomer)
	require.Equal(t, phoneNumber, *savedCustomer.GetPhoneNumber())

	savedCustomer, err = repository.FindById(ctx, aCustomer.Id)
	require.NoError(t, err)
	require.NotNil(t, savedCustomer)
	require.Equal(t, aCustomer, savedCustomer)
	require.Equal(t, phoneNumber, *savedCustomer.GetPhoneNumber())
}
