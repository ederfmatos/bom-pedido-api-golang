package repository

import (
	"bom-pedido-api/application/repository"
	customer2 "bom-pedido-api/domain/entity/customer"
	"context"
	"database/sql"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CustomerSqlRepository(t *testing.T) {
	database, err := sql.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	if err != nil {
		t.Error(err)
	}
	defer database.Close()
	_, err = database.Exec(`
		CREATE TABLE IF NOT EXISTS customers
		(
			id           VARCHAR(36)                NOT NULL PRIMARY KEY,
			name         VARCHAR(255)               NOT NULL,
			email        VARCHAR(255)               NOT NULL UNIQUE,
			phone_number VARCHAR(11)                UNIQUE,
			status       VARCHAR(11)				NOT NULL,
			created_at   TIMESTAMP                  NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	assert.NoError(t, err)
	sqlConnection := NewDefaultSqlConnection(database)
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
