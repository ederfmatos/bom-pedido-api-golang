package repository

import (
	"bom-pedido-api/internal/domain/entity/product"
	"context"
	"fmt"
)

type ProductMemoryRepository struct {
	products map[string]*product.Product
}

func NewProductMemoryRepository() *ProductMemoryRepository {
	return &ProductMemoryRepository{products: make(map[string]*product.Product)}
}

func (r *ProductMemoryRepository) Create(_ context.Context, product *product.Product) error {
	for _, p := range r.products {
		if p.Name == product.Name {
			return fmt.Errorf("duplicated product name")
		}
	}
	r.products[product.Id] = product
	return nil
}

func (r *ProductMemoryRepository) Update(_ context.Context, product *product.Product) error {
	r.products[product.Id] = product
	return nil
}

func (r *ProductMemoryRepository) FindById(_ context.Context, id string) (*product.Product, error) {
	return r.products[id], nil
}

func (r *ProductMemoryRepository) ExistsByNameAndTenantId(_ context.Context, name, tenantId string) (bool, error) {
	for _, aProduct := range r.products {
		if aProduct.Name == name && aProduct.TenantId == tenantId {
			return true, nil
		}
	}
	return false, nil
}

func (r *ProductMemoryRepository) FindAllById(_ context.Context, ids []string) (map[string]*product.Product, error) {
	products := make(map[string]*product.Product)
	for _, id := range ids {
		aProduct := r.products[id]
		if aProduct != nil {
			products[id] = aProduct
		}
	}
	return products, nil
}
