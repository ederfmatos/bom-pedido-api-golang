package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/pkg/mongo"
	"context"
	"fmt"
)

type ProductMongoRepository struct {
	collection mongo.Collection
}

func NewProductMongoRepository(database *mongo.Database) *ProductMongoRepository {
	return &ProductMongoRepository{collection: database.ForCollection("products")}
}

func (r *ProductMongoRepository) Create(ctx context.Context, product *entity.Product) error {
	return r.collection.InsertOne(ctx, product)
}

func (r *ProductMongoRepository) Update(ctx context.Context, product *entity.Product) error {
	return r.collection.UpdateByID(ctx, product.Id, product)
}

func (r *ProductMongoRepository) FindById(ctx context.Context, id string) (*entity.Product, error) {
	var product entity.Product
	err := r.collection.FindByID(ctx, id, &product)
	if err != nil || product.Id == "" {
		return nil, err
	}
	return &product, nil
}

func (r *ProductMongoRepository) ExistsByNameAndTenantId(ctx context.Context, name, tenantId string) (bool, error) {
	return r.collection.ExistsByTenantIdAnd(ctx, tenantId, "name", name)
}

func (r *ProductMongoRepository) FindAllById(ctx context.Context, ids []string) (map[string]*entity.Product, error) {
	products := make(map[string]*entity.Product)
	cursor, err := r.collection.FindAllByID(ctx, ids)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var product entity.Product
		if err = cursor.Decode(&product); err != nil {
			return nil, fmt.Errorf("decode product: %v", err)
		}
		products[product.Id] = &product
	}
	return products, nil
}
