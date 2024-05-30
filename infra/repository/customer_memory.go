package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity"
	"context"
)

type CustomerMemoryRepository struct {
	customers map[string]*entity.Customer
}

func NewCustomerMemoryRepository() repository.CustomerRepository {
	return &CustomerMemoryRepository{customers: make(map[string]*entity.Customer)}
}

func (repository *CustomerMemoryRepository) Create(_ context.Context, customer *entity.Customer) error {
	repository.customers[customer.Id] = customer
	return nil
}

func (repository *CustomerMemoryRepository) Update(_ context.Context, customer *entity.Customer) error {
	repository.customers[customer.Id] = customer
	return nil
}

func (repository *CustomerMemoryRepository) FindById(_ context.Context, id string) (*entity.Customer, error) {
	return repository.customers[id], nil
}

func (repository *CustomerMemoryRepository) FindByEmail(_ context.Context, email string) (*entity.Customer, error) {
	for _, customer := range repository.customers {
		if customer.Email.Value() == email {
			return customer, nil
		}
	}
	return nil, nil
}
