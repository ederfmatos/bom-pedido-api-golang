package repository

import (
	"bom-pedido-api/internal/domain/entity/customer"
	"bom-pedido-api/pkg/mongo"
	"context"
)

type CustomerMongoRepository struct {
	collection *mongo.Collection
}

func NewCustomerMongoRepository(mongoDatabase *mongo.Database) *CustomerMongoRepository {
	return &CustomerMongoRepository{collection: mongoDatabase.ForCollection("customers")}
}

func (r *CustomerMongoRepository) Create(ctx context.Context, customer *customer.Customer) error {
	return r.collection.InsertOne(ctx, customer)
}

func (r *CustomerMongoRepository) Update(ctx context.Context, customer *customer.Customer) error {
	return r.collection.UpdateByID(ctx, customer.Id, customer)
}

func (r *CustomerMongoRepository) FindById(ctx context.Context, id string) (*customer.Customer, error) {
	var aCustomer customer.Customer
	err := r.collection.FindByID(ctx, id, &aCustomer)
	if err != nil || aCustomer.Id == "" {
		return nil, err
	}
	return &aCustomer, nil
}

func (r *CustomerMongoRepository) FindByEmail(ctx context.Context, email, tenantId string) (*customer.Customer, error) {
	var aCustomer customer.Customer
	err := r.collection.FindByTenantIdAnd(ctx, tenantId, "email", email, &aCustomer)
	if err != nil || aCustomer.Id == "" {
		return nil, err
	}
	return &aCustomer, nil
}
