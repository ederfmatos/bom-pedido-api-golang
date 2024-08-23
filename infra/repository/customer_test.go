package repository

import (
	"bom-pedido-api/application/repository"
	customer "bom-pedido-api/domain/entity/customer"
	"bom-pedido-api/infra/test"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
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
	ctx := context.TODO()

	aCustomer, err := customer.New(faker.Name(), faker.Email(), faker.Word())
	assert.NoError(t, err)

	savedCustomer, err := repository.FindByEmail(ctx, aCustomer.GetEmail(), aCustomer.TenantId)
	assert.NoError(t, err)
	assert.Nil(t, savedCustomer)

	savedCustomer, err = repository.FindById(ctx, aCustomer.Id)
	assert.NoError(t, err)
	assert.Nil(t, savedCustomer)

	err = repository.Create(ctx, aCustomer)
	assert.NoError(t, err)

	savedCustomer, err = repository.FindByEmail(ctx, aCustomer.GetEmail(), aCustomer.TenantId)
	assert.NoError(t, err)
	assert.NotNil(t, savedCustomer)
	assert.Equal(t, aCustomer, savedCustomer)

	savedCustomer, err = repository.FindById(ctx, aCustomer.Id)
	assert.NoError(t, err)
	assert.NotNil(t, savedCustomer)
	assert.Equal(t, aCustomer, savedCustomer)

	phoneNumber := "11999999999"
	err = aCustomer.SetPhoneNumber(phoneNumber)
	assert.NoError(t, err)

	err = repository.Update(ctx, aCustomer)
	assert.NoError(t, err)

	savedCustomer, err = repository.FindByEmail(ctx, aCustomer.GetEmail(), aCustomer.TenantId)
	assert.NoError(t, err)
	assert.NotNil(t, savedCustomer)
	assert.Equal(t, aCustomer, savedCustomer)
	assert.Equal(t, phoneNumber, *savedCustomer.GetPhoneNumber())

	savedCustomer, err = repository.FindById(ctx, aCustomer.Id)
	assert.NoError(t, err)
	assert.NotNil(t, savedCustomer)
	assert.Equal(t, aCustomer, savedCustomer)
	assert.Equal(t, phoneNumber, *savedCustomer.GetPhoneNumber())
}
