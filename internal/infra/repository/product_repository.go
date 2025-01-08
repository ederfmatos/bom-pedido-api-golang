package repository

import (
	"bom-pedido-api/internal/domain/entity/product"
	"bom-pedido-api/pkg/mongo"
	"context"
	"fmt"
)

type ProductMongoRepository struct {
	collection *mongo.Collection
}

func NewProductMongoRepository(database *mongo.Database) *ProductMongoRepository {
	return &ProductMongoRepository{collection: database.ForCollection("products")}
}

func (r *ProductMongoRepository) Create(ctx context.Context, product *product.Product) error {
	return r.collection.InsertOne(ctx, product)
}

func (r *ProductMongoRepository) Update(ctx context.Context, product *product.Product) error {
	return r.collection.UpdateByID(ctx, product.Id, product)
}

func (r *ProductMongoRepository) FindById(ctx context.Context, id string) (*product.Product, error) {
	var aProduct product.Product
	err := r.collection.FindByID(ctx, id, &aProduct)
	if err != nil || aProduct.Id == "" {
		return nil, err
	}
	return &aProduct, nil
}

func (r *ProductMongoRepository) ExistsByNameAndTenantId(ctx context.Context, name, tenantId string) (bool, error) {
	return r.collection.ExistsByTenantIdAnd(ctx, tenantId, "name", name)
}

func (r *ProductMongoRepository) FindAllById(ctx context.Context, ids []string) (map[string]*product.Product, error) {
	products := make(map[string]*product.Product)
	cursor, err := r.collection.FindAllByID(ctx, ids)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var aProduct product.Product
		if err = cursor.Decode(&aProduct); err != nil {
			return nil, fmt.Errorf("decode product: %v", err)
		}
		products[aProduct.Id] = &aProduct
	}
	return products, nil
}
