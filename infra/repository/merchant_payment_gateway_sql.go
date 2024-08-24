package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/merchant"
	"context"
)

const (
	sqlFindMerchantPaymentGatewayConfig = "SELECT credentials FROM merchant_payment_gateway_configs WHERE merchant_id = $1 AND gateway = $2 LIMIT 1"
)

type DefaultMerchantPaymentGatewayConfigRepository struct {
	SqlConnection
}

func NewDefaultMerchantPaymentGatewayConfigRepository(sqlConnection SqlConnection) repository.MerchantPaymentGatewayConfigRepository {
	return &DefaultMerchantPaymentGatewayConfigRepository{SqlConnection: sqlConnection}
}

func (r *DefaultMerchantPaymentGatewayConfigRepository) FindByMerchantAndGateway(ctx context.Context, merchantId, gateway string) (*merchant.PaymentGatewayConfig, error) {
	var accessToken string
	found, err := r.Sql(sqlFindMerchantPaymentGatewayConfig).Values(merchantId, gateway).FindOne(ctx, &accessToken)
	if err != nil || !found {
		return nil, err
	}
	return &merchant.PaymentGatewayConfig{
		MerchantID:  merchantId,
		AccessToken: accessToken,
		Gateway:     gateway,
	}, nil
}
