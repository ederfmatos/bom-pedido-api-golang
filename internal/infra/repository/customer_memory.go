package repository

import (
	"bom-pedido-api/internal/domain/entity/customer"
	"context"
)

type CustomerMemoryRepository struct {
	customers map[string]*customer.Customer
}

func NewCustomerMemoryRepository() *CustomerMemoryRepository {
	return &CustomerMemoryRepository{customers: make(map[string]*customer.Customer)}
}

func (r *CustomerMemoryRepository) Create(_ context.Context, customer *customer.Customer) error {
	r.customers[customer.Id] = customer
	return nil
}

func (r *CustomerMemoryRepository) Update(_ context.Context, customer *customer.Customer) error {
	r.customers[customer.Id] = customer
	return nil
}

func (r *CustomerMemoryRepository) FindById(_ context.Context, id string) (*customer.Customer, error) {
	return r.customers[id], nil
}

func (r *CustomerMemoryRepository) FindByEmail(_ context.Context, email, tenantId string) (*customer.Customer, error) {
	for _, c := range r.customers {
		if c.GetEmail() == email && c.TenantId == tenantId {
			return c, nil
		}
	}
	return nil, nil
}
