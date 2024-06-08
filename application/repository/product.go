package repository

import (
	"bom-pedido-api/domain/entity"
	"context"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	Update(ctx context.Context, product *entity.Product) error
	FindById(ctx context.Context, id string) (*entity.Product, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
}
