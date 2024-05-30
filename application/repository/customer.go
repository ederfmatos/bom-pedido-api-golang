package repository

import (
	"bom-pedido-api/domain/entity"
	"context"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *entity.Customer) error
	Update(ctx context.Context, customer *entity.Customer) error
	FindById(ctx context.Context, id string) (*entity.Customer, error)
	FindByEmail(ctx context.Context, email string) (*entity.Customer, error)
}
