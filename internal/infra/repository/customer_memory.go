package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"context"
)

type CustomerMemoryRepository struct {
	customers map[string]*entity.Customer
}

func NewCustomerMemoryRepository() *CustomerMemoryRepository {
	return &CustomerMemoryRepository{customers: make(map[string]*entity.Customer)}
}

func (r *CustomerMemoryRepository) Create(_ context.Context, customer *entity.Customer) error {
	r.customers[customer.Id] = customer
	return nil
}

func (r *CustomerMemoryRepository) Update(_ context.Context, customer *entity.Customer) error {
	r.customers[customer.Id] = customer
	return nil
}

func (r *CustomerMemoryRepository) FindById(_ context.Context, id string) (*entity.Customer, error) {
	return r.customers[id], nil
}

func (r *CustomerMemoryRepository) FindByEmail(_ context.Context, email, tenantId string) (*entity.Customer, error) {
	for _, c := range r.customers {
		if c.GetEmail() == email && c.TenantId == tenantId {
			return c, nil
		}
	}
	return nil, nil
}
