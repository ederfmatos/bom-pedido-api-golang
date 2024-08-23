package repository

import (
	"bom-pedido-api/domain/entity/customer"
	"context"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *customer.Customer) error
	Update(ctx context.Context, customer *customer.Customer) error
	FindById(ctx context.Context, id string) (*customer.Customer, error)
	FindByEmail(ctx context.Context, email, tenantId string) (*customer.Customer, error)
}
