package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"context"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *entity.Customer) error
	Update(ctx context.Context, customer *entity.Customer) error
	FindById(ctx context.Context, id string) (*entity.Customer, error)
	FindByEmail(ctx context.Context, email, tenantId string) (*entity.Customer, error)
}
