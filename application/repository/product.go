package repository

import (
	"bom-pedido-api/domain/entity/product"
	"context"
)

type ProductRepository interface {
	Create(ctx context.Context, product *product.Product) error
	Update(ctx context.Context, product *product.Product) error
	FindById(ctx context.Context, id string) (*product.Product, error)
	ExistsByNameAndTenantId(ctx context.Context, name, tenantId string) (bool, error)
	FindAllById(ctx context.Context, ids []string) (map[string]*product.Product, error)
}
