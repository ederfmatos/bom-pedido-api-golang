package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/product"
	"context"
	"strconv"
)

const (
	sqlCreateProduct       = "INSERT INTO products (id, name, description, price, status, category_id, tenant_id) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	sqlUpdateProduct       = "UPDATE products SET name = $1, description = $2, price = $3, status = $4 WHERE id = $5"
	sqlFindProductById     = "SELECT id, name, description, price, status, category_id, tenant_id FROM products WHERE id = $1"
	sqlExistsProductByName = "SELECT 1 FROM products WHERE name = $1 AND tenant_id = $2 LIMIT 1"
)

type DefaultProductRepository struct {
	SqlConnection
}

func NewDefaultProductRepository(sqlConnection SqlConnection) repository.ProductRepository {
	return &DefaultProductRepository{SqlConnection: sqlConnection}
}

func (repository *DefaultProductRepository) Create(ctx context.Context, product *product.Product) error {
	return repository.Sql(sqlCreateProduct).
		Values(product.Id, product.Name, product.Description, product.Price, product.Status, product.CategoryId, product.TenantId).
		Update(ctx)
}

func (repository *DefaultProductRepository) Update(ctx context.Context, product *product.Product) error {
	return repository.Sql(sqlUpdateProduct).
		Values(product.Name, product.Description, product.Price, product.Status, product.Id).
		Update(ctx)
}

func (repository *DefaultProductRepository) FindById(ctx context.Context, id string) (*product.Product, error) {
	var name, description, status, tenantId, categoryId string
	var price float64
	found, err := repository.Sql(sqlFindProductById).
		Values(id).
		FindOne(ctx, &id, &name, &description, &price, &status, &categoryId, &tenantId)
	if err != nil || !found {
		return nil, err
	}
	return product.Restore(id, name, description, price, status, categoryId, tenantId)
}

func (repository *DefaultProductRepository) ExistsByNameAndTenantId(ctx context.Context, name, tenantId string) (bool, error) {
	return repository.Sql(sqlExistsProductByName).
		Values(name, tenantId).
		Exists(ctx)
}

func (repository *DefaultProductRepository) FindAllById(ctx context.Context, ids []string) (map[string]*product.Product, error) {
	products := make(map[string]*product.Product)
	args := make([]interface{}, len(ids))
	in := ""
	for i, id := range ids {
		args[i] = id
		if i > 0 {
			in += ","
		}
		in += "$" + strconv.Itoa(i+1)
	}
	var sqlListProducts = `select id, name, description, price, status, tenant_id, category_id from products WHERE id IN (` + in + `)`
	err := repository.Sql(sqlListProducts).
		Values(args...).
		List(ctx, func(getValues func(dest ...any) error) error {
			var id, name, description, status, categoryId, tenantId string
			var price float64
			err := getValues(&id, &name, &description, &price, &status, &categoryId, &tenantId)
			if err != nil {
				return err
			}
			aProduct, err := product.Restore(id, name, description, price, status, categoryId, tenantId)
			if err != nil {
				return err
			}
			products[id] = aProduct
			return nil
		})
	return products, err
}
