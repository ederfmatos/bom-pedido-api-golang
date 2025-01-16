package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"context"
	"fmt"
)

type ProductMemoryRepository struct {
	products map[string]*entity.Product
}

func NewProductMemoryRepository() *ProductMemoryRepository {
	return &ProductMemoryRepository{products: make(map[string]*entity.Product)}
}

func (r *ProductMemoryRepository) Create(_ context.Context, product *entity.Product) error {
	for _, p := range r.products {
		if p.Name == product.Name {
			return fmt.Errorf("duplicated product name")
		}
	}
	r.products[product.Id] = product
	return nil
}

func (r *ProductMemoryRepository) Update(_ context.Context, product *entity.Product) error {
	r.products[product.Id] = product
	return nil
}

func (r *ProductMemoryRepository) FindById(_ context.Context, id string) (*entity.Product, error) {
	return r.products[id], nil
}

func (r *ProductMemoryRepository) ExistsByNameAndTenantId(_ context.Context, name, tenantId string) (bool, error) {
	for _, product := range r.products {
		if product.Name == name && product.TenantId == tenantId {
			return true, nil
		}
	}
	return false, nil
}

func (r *ProductMemoryRepository) FindAllById(_ context.Context, ids []string) (map[string]*entity.Product, error) {
	products := make(map[string]*entity.Product)
	for _, id := range ids {
		product := r.products[id]
		if product != nil {
			products[id] = product
		}
	}
	return products, nil
}
