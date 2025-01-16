package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"context"
)

type CategoriesMemoryRepository struct {
	categories map[string]*entity.Category
}

func NewCategoriesMemoryRepository() *CategoriesMemoryRepository {
	return &CategoriesMemoryRepository{categories: make(map[string]*entity.Category)}
}

func (r *CategoriesMemoryRepository) Create(_ context.Context, product *entity.Category) error {
	r.categories[product.Id] = product
	return nil
}

func (r *CategoriesMemoryRepository) Update(_ context.Context, product *entity.Category) error {
	r.categories[product.Id] = product
	return nil
}

func (r *CategoriesMemoryRepository) ExistsById(_ context.Context, id string) (bool, error) {
	_, found := r.categories[id]
	return found, nil
}

func (r *CategoriesMemoryRepository) ExistsByNameAndTenantId(_ context.Context, name, tenantId string) (bool, error) {
	for _, category := range r.categories {
		if category.Name == name && category.TenantId == tenantId {
			return true, nil
		}
	}
	return false, nil
}
