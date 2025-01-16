package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"context"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	Update(ctx context.Context, product *entity.Product) error
	FindById(ctx context.Context, id string) (*entity.Product, error)
	ExistsByNameAndTenantId(ctx context.Context, name, tenantId string) (bool, error)
	FindAllById(ctx context.Context, ids []string) (map[string]*entity.Product, error)
}

type ProductCategoryRepository interface {
	Create(context.Context, *entity.Category) error
	ExistsByNameAndTenantId(ctx context.Context, name, tenantId string) (bool, error)
	ExistsById(ctx context.Context, id string) (bool, error)
}
