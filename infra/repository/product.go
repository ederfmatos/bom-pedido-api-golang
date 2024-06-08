package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity"
	"context"
)

type DefaultProductRepository struct {
	SqlConnection
}

func NewDefaultProductRepository(sqlConnection SqlConnection) repository.ProductRepository {
	return &DefaultProductRepository{SqlConnection: sqlConnection}
}

func (repository *DefaultProductRepository) Create(ctx context.Context, product *entity.Product) error {
	return repository.Sql("INSERT INTO products (id, name, description, price, status) VALUES (?, ?, ?, ?, ?)").
		Values(product.ID, product.Name, product.Description, product.Price, product.Status).
		Update(ctx)
}

func (repository *DefaultProductRepository) Update(ctx context.Context, product *entity.Product) error {
	return repository.Sql("UPDATE products SET name = ?, description = ?, price = ?, status = ? WHERE id = ?").
		Values(product.Name, product.Description, product.Price, product.Status, product.ID).
		Update(ctx)
}

func (repository *DefaultProductRepository) FindById(ctx context.Context, id string) (*entity.Product, error) {
	var name string
	var description string
	var price float64
	var status string
	found, err := repository.Sql("SELECT id, name, description, price, status FROM products WHERE id = ?").
		Values(id).
		FindOne(ctx, &id, &name, &description, &price, &status)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return entity.RestoreProduct(id, name, description, price, status)
}

func (repository *DefaultProductRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	return repository.Sql("SELECT 1 FROM products WHERE name = ? LIMIT 1").
		Values(name).
		Exists(ctx)
}
