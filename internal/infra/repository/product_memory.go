package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/product"
	"context"
	"fmt"
)

type ProductMemoryRepository struct {
	products map[string]*product.Product
}

func NewProductMemoryRepository() repository.ProductRepository {
	return &ProductMemoryRepository{products: make(map[string]*product.Product)}
}

func (repository *ProductMemoryRepository) Create(_ context.Context, product *product.Product) error {
	for _, p := range repository.products {
		if p.Name == product.Name {
			return fmt.Errorf("duplicated product name")
		}
	}
	repository.products[product.Id] = product
	return nil
}

func (repository *ProductMemoryRepository) Update(_ context.Context, product *product.Product) error {
	repository.products[product.Id] = product
	return nil
}

func (repository *ProductMemoryRepository) FindById(_ context.Context, id string) (*product.Product, error) {
	return repository.products[id], nil
}

func (repository *ProductMemoryRepository) ExistsByNameAndTenantId(_ context.Context, name, tenantId string) (bool, error) {
	for _, aProduct := range repository.products {
		if aProduct.Name == name && aProduct.TenantId == tenantId {
			return true, nil
		}
	}
	return false, nil
}

func (repository *ProductMemoryRepository) FindAllById(_ context.Context, ids []string) (map[string]*product.Product, error) {
	products := make(map[string]*product.Product)
	for _, id := range ids {
		aProduct := repository.products[id]
		if aProduct != nil {
			products[id] = aProduct
		}
	}
	return products, nil
}
