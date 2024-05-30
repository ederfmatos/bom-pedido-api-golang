package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity"
	"context"
)

type DefaultCustomerRepository struct {
	SqlConnection
}

func NewDefaultCustomerRepository(sqlConnection SqlConnection) repository.CustomerRepository {
	return &DefaultCustomerRepository{SqlConnection: sqlConnection}
}

func (repository *DefaultCustomerRepository) Create(ctx context.Context, customer *entity.Customer) error {
	return repository.Sql("INSERT INTO customers (id, name, email, status) VALUES (?, ?, ?, ?)").
		Values(customer.Id, customer.Name, customer.Email, customer.Status).
		Update(ctx)
}

func (repository *DefaultCustomerRepository) Update(ctx context.Context, customer *entity.Customer) error {
	return repository.Sql("UPDATE customers SET name = ?, email = ?, status = ?, phone_number = ? WHERE id = ?").
		Values(customer.Name, customer.Email, customer.Status, customer.PhoneNumber, customer.Id).
		Update(ctx)
}

func (repository *DefaultCustomerRepository) FindById(ctx context.Context, id string) (*entity.Customer, error) {
	var customer entity.Customer
	err := repository.Sql("SELECT id, name, email, phone_number, status FROM customers WHERE id = ?").
		Values(id).
		FindOne(ctx, &customer.Id, &customer.Name, &customer.Email, &customer.PhoneNumber, &customer.Status)
	return &customer, err
}

func (repository *DefaultCustomerRepository) FindByEmail(ctx context.Context, email string) (*entity.Customer, error) {
	var customer entity.Customer
	err := repository.Sql("SELECT id, name, email, phone_number, status FROM customers WHERE email = ?").
		Values(email).
		FindOne(ctx, &customer.Id, &customer.Name, &customer.Email, &customer.PhoneNumber, &customer.Status)
	return &customer, err
}
