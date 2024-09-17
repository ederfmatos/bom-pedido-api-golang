package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/product"
	"context"
)

const (
	sqlCreateProductCategory       = "INSERT INTO product_categories (id, name, description, tenant_id) VALUES ($1, $2, $3, $4)"
	sqlExistsProductCategoryByName = "SELECT 1 FROM product_categories WHERE name = $1 AND tenant_id = $2 LIMIT 1"
	sqlExistsProductCategoryById   = "SELECT 1 FROM product_categories WHERE id = $1 LIMIT 1"
)

type DefaultProductCategoryRepository struct {
	SqlConnection
}

func NewDefaultProductCategoryRepository(sqlConnection SqlConnection) repository.ProductCategoryRepository {
	return &DefaultProductCategoryRepository{SqlConnection: sqlConnection}
}

func (repository *DefaultProductCategoryRepository) Create(ctx context.Context, product *product.Category) error {
	return repository.Sql(sqlCreateProductCategory).
		Values(product.Id, product.Name, product.Description, product.TenantId).
		Update(ctx)
}

func (repository *DefaultProductCategoryRepository) ExistsById(ctx context.Context, id string) (bool, error) {
	return repository.Sql(sqlExistsProductCategoryById).
		Values(id).
		Exists(ctx)
}

func (repository *DefaultProductCategoryRepository) ExistsByNameAndTenantId(ctx context.Context, name, tenantId string) (bool, error) {
	return repository.Sql(sqlExistsProductCategoryByName).
		Values(name, tenantId).
		Exists(ctx)
}
