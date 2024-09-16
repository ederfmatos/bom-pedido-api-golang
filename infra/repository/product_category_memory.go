package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/product"
	"context"
)

type ProductCategoryMemoryRepository struct {
	categories map[string]*product.Category
}

func NewProductCategoryMemoryRepository() repository.ProductCategoryRepository {
	return &ProductCategoryMemoryRepository{categories: make(map[string]*product.Category)}
}

func (repository *ProductCategoryMemoryRepository) Create(_ context.Context, product *product.Category) error {
	repository.categories[product.Id] = product
	return nil
}

func (repository *ProductCategoryMemoryRepository) Update(_ context.Context, product *product.Category) error {
	repository.categories[product.Id] = product
	return nil
}

func (repository *ProductCategoryMemoryRepository) ExistsByNameAndTenantId(_ context.Context, name, tenantId string) (bool, error) {
	for _, aProductCategory := range repository.categories {
		if aProductCategory.Name == name && aProductCategory.TenantId == tenantId {
			return true, nil
		}
	}
	return false, nil
}
