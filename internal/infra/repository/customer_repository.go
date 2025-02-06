package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/pkg/mongo"
	"context"
)

type CustomerMongoRepository struct {
	collection mongo.Collection
}

func NewCustomerMongoRepository(mongoDatabase *mongo.Database) *CustomerMongoRepository {
	return &CustomerMongoRepository{collection: mongoDatabase.ForCollection("customers")}
}

func (r *CustomerMongoRepository) Create(ctx context.Context, customer *entity.Customer) error {
	return r.collection.InsertOne(ctx, customer)
}

func (r *CustomerMongoRepository) Update(ctx context.Context, customer *entity.Customer) error {
	return r.collection.UpdateByID(ctx, customer.Id, customer)
}

func (r *CustomerMongoRepository) FindById(ctx context.Context, id string) (*entity.Customer, error) {
	var customer entity.Customer
	err := r.collection.FindByID(ctx, id, &customer)
	if err != nil || customer.Id == "" {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerMongoRepository) FindByEmail(ctx context.Context, email, tenantId string) (*entity.Customer, error) {
	var customer entity.Customer
	err := r.collection.FindByTenantIdAnd(ctx, tenantId, "email", email, &customer)
	if err != nil || customer.Id == "" {
		return nil, err
	}
	return &customer, nil
}
