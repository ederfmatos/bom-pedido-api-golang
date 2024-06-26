package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity"
	"context"
)

type ProductMemoryRepository struct {
	products map[string]*entity.Product
}

func NewProductMemoryRepository() repository.ProductRepository {
	return &ProductMemoryRepository{products: make(map[string]*entity.Product)}
}

func (repository *ProductMemoryRepository) Create(_ context.Context, product *entity.Product) error {
	repository.products[product.ID] = product
	return nil
}

func (repository *ProductMemoryRepository) Update(_ context.Context, product *entity.Product) error {
	repository.products[product.ID] = product
	return nil
}

func (repository *ProductMemoryRepository) FindById(_ context.Context, id string) (*entity.Product, error) {
	return repository.products[id], nil
}

func (repository *ProductMemoryRepository) ExistsByName(_ context.Context, name string) (bool, error) {
	for _, product := range repository.products {
		if product.Name == name {
			return true, nil
		}
	}
	return false, nil
}
