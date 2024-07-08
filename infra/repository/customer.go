package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/customer"
	"context"
)

const (
	sqlCreateCustomer      = "INSERT INTO customers (id, name, email, status, phone_number) VALUES (?, ?, ?, ?, ?)"
	sqlUpdateCustomer      = "UPDATE customers SET name = ?, email = ?, status = ?, phone_number = ? WHERE id = ?"
	sqlFindCustomerById    = "SELECT id, name, email, phone_number, status FROM customers WHERE id = ?"
	sqlFindCustomerByEmail = "SELECT id, name, email, phone_number, status FROM customers WHERE email = ?"
)

type DefaultCustomerRepository struct {
	SqlConnection
}

func NewDefaultCustomerRepository(sqlConnection SqlConnection) repository.CustomerRepository {
	return &DefaultCustomerRepository{SqlConnection: sqlConnection}
}

func (repository *DefaultCustomerRepository) Create(ctx context.Context, customer *customer.Customer) error {
	return repository.Sql(sqlCreateCustomer).
		Values(customer.Id, customer.Name, customer.GetEmail(), customer.Status, customer.GetPhoneNumber()).
		Update(ctx)
}

func (repository *DefaultCustomerRepository) Update(ctx context.Context, customer *customer.Customer) error {
	return repository.Sql(sqlUpdateCustomer).
		Values(customer.Name, customer.GetEmail(), customer.Status, customer.GetPhoneNumber(), customer.Id).
		Update(ctx)
}

func (repository *DefaultCustomerRepository) FindById(ctx context.Context, id string) (*customer.Customer, error) {
	var email string
	var name string
	var status string
	var phoneNumber string
	found, err := repository.Sql(sqlFindCustomerById).
		Values(id).
		FindOne(ctx, &id, &name, &email, &phoneNumber, &status)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return customer.Restore(id, name, email, phoneNumber, status)
}

func (repository *DefaultCustomerRepository) FindByEmail(ctx context.Context, email string) (*customer.Customer, error) {
	var id string
	var name string
	var status string
	var phoneNumber string
	found, err := repository.Sql(sqlFindCustomerByEmail).
		Values(email).
		FindOne(ctx, &id, &name, &email, &phoneNumber, &status)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return customer.Restore(id, name, email, phoneNumber, status)
}
