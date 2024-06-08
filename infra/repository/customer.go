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
	return repository.Sql("INSERT INTO customers (id, name, email, status, phone_number) VALUES (?, ?, ?, ?, ?)").
		Values(customer.Id, customer.Name, customer.GetEmail(), customer.Status, customer.GetPhoneNumber()).
		Update(ctx)
}

func (repository *DefaultCustomerRepository) Update(ctx context.Context, customer *entity.Customer) error {
	return repository.Sql("UPDATE customers SET name = ?, email = ?, status = ?, phone_number = ? WHERE id = ?").
		Values(customer.Name, customer.GetEmail(), customer.Status, customer.GetPhoneNumber(), customer.Id).
		Update(ctx)
}

func (repository *DefaultCustomerRepository) FindById(ctx context.Context, id string) (*entity.Customer, error) {
	var email string
	var name string
	var status string
	var phoneNumber string
	found, err := repository.Sql("SELECT id, name, phone_number, status FROM customers WHERE id = ?").
		Values(id).
		FindOne(ctx, &id, &name, &email, &phoneNumber, &status)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return entity.RestoreCustomer(id, name, email, phoneNumber, status)
}

func (repository *DefaultCustomerRepository) FindByEmail(ctx context.Context, email string) (*entity.Customer, error) {
	var id string
	var name string
	var status string
	var phoneNumber string
	found, err := repository.Sql("SELECT id, name, phone_number, status FROM customers WHERE email = ?").
		Values(email).
		FindOne(ctx, &id, &name, &email, &phoneNumber, &status)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return entity.RestoreCustomer(id, name, email, phoneNumber, status)
}
