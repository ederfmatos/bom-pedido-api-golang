package repository

import (
	"bom-pedido-api/internal/domain/entity/merchant"
	"bom-pedido-api/pkg/mongo"
	"context"
)

type MerchantPaymentGatewayConfigMongoRepository struct {
	collection *mongo.Collection
}

func NewMerchantPaymentGatewayConfigMongoRepository(database *mongo.Database) *MerchantPaymentGatewayConfigMongoRepository {
	return &MerchantPaymentGatewayConfigMongoRepository{collection: database.ForCollection("merchant_payment_gateway_config")}
}

func (r *MerchantPaymentGatewayConfigMongoRepository) Create(ctx context.Context, config *merchant.PaymentGatewayConfig) error {
	return r.collection.InsertOne(ctx, config)
}

func (r *MerchantPaymentGatewayConfigMongoRepository) FindByMerchantAndGateway(ctx context.Context, merchantId, gateway string) (*merchant.PaymentGatewayConfig, error) {
	var config merchant.PaymentGatewayConfig
	err := r.collection.FindByValues(ctx, map[string]interface{}{"merchantId": merchantId, "paymentGateway": gateway}, &config)
	if err != nil || config.MerchantID == "" {
		return nil, err
	}
	return &config, nil
}

func (r *MerchantPaymentGatewayConfigMongoRepository) FindByMerchant(ctx context.Context, merchantId string) (*merchant.PaymentGatewayConfig, error) {
	var config merchant.PaymentGatewayConfig
	err := r.collection.FindBy(ctx, "merchantId", merchantId, &config)
	if err != nil || config.MerchantID == "" {
		return nil, err
	}
	return &config, nil
}
