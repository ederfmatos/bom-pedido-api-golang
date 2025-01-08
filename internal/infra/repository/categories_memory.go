package repository

import (
	"bom-pedido-api/internal/domain/entity/product"
	"context"
)

type CategoriesMemoryRepository struct {
	categories map[string]*product.Category
}

func NewCategoriesMemoryRepository() *CategoriesMemoryRepository {
	return &CategoriesMemoryRepository{categories: make(map[string]*product.Category)}
}

func (r *CategoriesMemoryRepository) Create(_ context.Context, product *product.Category) error {
	r.categories[product.Id] = product
	return nil
}

func (r *CategoriesMemoryRepository) Update(_ context.Context, product *product.Category) error {
	r.categories[product.Id] = product
	return nil
}

func (r *CategoriesMemoryRepository) ExistsById(_ context.Context, id string) (bool, error) {
	_, found := r.categories[id]
	return found, nil
}

func (r *CategoriesMemoryRepository) ExistsByNameAndTenantId(_ context.Context, name, tenantId string) (bool, error) {
	for _, aProductCategory := range r.categories {
		if aProductCategory.Name == name && aProductCategory.TenantId == tenantId {
			return true, nil
		}
	}
	return false, nil
}
