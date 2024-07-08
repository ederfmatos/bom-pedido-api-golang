package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/product"
	"context"
	"strings"
)

const (
	sqlCreateProduct       = "INSERT INTO products (id, name, description, price, status) VALUES (?, ?, ?, ?, ?)"
	sqlUpdateProduct       = "UPDATE products SET name = ?, description = ?, price = ?, status = ? WHERE id = ?"
	sqlFindProductById     = "SELECT id, name, description, price, status FROM products WHERE id = ?"
	sqlExistsProductByName = "SELECT 1 FROM products WHERE name = ? LIMIT 1"
	sqlListProducts        = "select id, name, description, price, status from products WHERE id IN (?)"
)

type DefaultProductRepository struct {
	SqlConnection
}

func NewDefaultProductRepository(sqlConnection SqlConnection) repository.ProductRepository {
	return &DefaultProductRepository{SqlConnection: sqlConnection}
}

func (repository *DefaultProductRepository) Create(ctx context.Context, product *product.Product) error {
	return repository.Sql(sqlCreateProduct).
		Values(product.Id, product.Name, product.Description, product.Price, product.Status).
		Update(ctx)
}

func (repository *DefaultProductRepository) Update(ctx context.Context, product *product.Product) error {
	return repository.Sql(sqlUpdateProduct).
		Values(product.Name, product.Description, product.Price, product.Status, product.Id).
		Update(ctx)
}

func (repository *DefaultProductRepository) FindById(ctx context.Context, id string) (*product.Product, error) {
	var name, description, status string
	var price float64
	found, err := repository.Sql(sqlFindProductById).
		Values(id).
		FindOne(ctx, &id, &name, &description, &price, &status)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return product.Restore(id, name, description, price, status)
}

func (repository *DefaultProductRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	return repository.Sql(sqlExistsProductByName).
		Values(name).
		Exists(ctx)
}

func (repository *DefaultProductRepository) FindAllById(ctx context.Context, ids []string) (map[string]*product.Product, error) {
	products := make(map[string]*product.Product)
	err := repository.Sql(sqlListProducts).
		Values(strings.Join(ids, "','")).
		List(ctx, func(getValues func(dest ...any) error) error {
			var id, name, description, status string
			var price float64
			err := getValues(&id, &name, &description, &price, &status)
			if err != nil {
				return err
			}
			product, err := product.Restore(id, name, description, price, status)
			if err != nil {
				return err
			}
			products[id] = product
			return nil
		})
	return products, err
}
