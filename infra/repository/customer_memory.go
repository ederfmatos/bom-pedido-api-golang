package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/customer"
	"context"
)

type CustomerMemoryRepository struct {
	customers map[string]*customer.Customer
}

func NewCustomerMemoryRepository() repository.CustomerRepository {
	return &CustomerMemoryRepository{customers: make(map[string]*customer.Customer)}
}

func (repository *CustomerMemoryRepository) Create(_ context.Context, customer *customer.Customer) error {
	repository.customers[customer.Id] = customer
	return nil
}

func (repository *CustomerMemoryRepository) Update(_ context.Context, customer *customer.Customer) error {
	repository.customers[customer.Id] = customer
	return nil
}

func (repository *CustomerMemoryRepository) FindById(_ context.Context, id string) (*customer.Customer, error) {
	return repository.customers[id], nil
}

func (repository *CustomerMemoryRepository) FindByEmail(_ context.Context, email, tenantId string) (*customer.Customer, error) {
	for _, c := range repository.customers {
		if *c.GetEmail() == email && c.TenantId == tenantId {
			return c, nil
		}
	}
	return nil, nil
}
