package repository

import (
	"bom-pedido-api/application/repository"
	customer2 "bom-pedido-api/domain/entity/customer"
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

	customer, err := customer2.New(faker.Name(), faker.Email())
	assert.NoError(t, err)

	savedCustomer, err := repository.FindByEmail(ctx, *customer.GetEmail())
	assert.NoError(t, err)
	assert.Nil(t, savedCustomer)

	savedCustomer, err = repository.FindById(ctx, customer.Id)
	assert.NoError(t, err)
	assert.Nil(t, savedCustomer)

	err = repository.Create(ctx, customer)
	assert.NoError(t, err)

	savedCustomer, err = repository.FindByEmail(ctx, *customer.GetEmail())
	assert.NoError(t, err)
	assert.NotNil(t, savedCustomer)
	assert.Equal(t, customer, savedCustomer)

	savedCustomer, err = repository.FindById(ctx, customer.Id)
	assert.NoError(t, err)
	assert.NotNil(t, savedCustomer)
	assert.Equal(t, customer, savedCustomer)

	phoneNumber := "11999999999"
	err = customer.SetPhoneNumber(phoneNumber)
	assert.NoError(t, err)

	err = repository.Update(ctx, customer)
	assert.NoError(t, err)

	savedCustomer, err = repository.FindByEmail(ctx, *customer.GetEmail())
	assert.NoError(t, err)
	assert.NotNil(t, savedCustomer)
	assert.Equal(t, customer, savedCustomer)
	assert.Equal(t, phoneNumber, *savedCustomer.GetPhoneNumber())

	savedCustomer, err = repository.FindById(ctx, customer.Id)
	assert.NoError(t, err)
	assert.NotNil(t, savedCustomer)
	assert.Equal(t, customer, savedCustomer)
	assert.Equal(t, phoneNumber, *savedCustomer.GetPhoneNumber())
}
