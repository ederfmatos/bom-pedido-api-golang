package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/merchant"
	"context"
)

const (
	sqlFindMerchantPaymentGatewayConfigByMerchant           = "SELECT credentials, gateway FROM merchant_payment_gateway_configs WHERE merchant_id = $1 AND status = 'ACTIVE' ORDER BY priority LIMIT 1"
	sqlFindMerchantPaymentGatewayConfigByMerchantAndGateway = "SELECT credentials FROM merchant_payment_gateway_configs WHERE merchant_id = $1 AND gateway = $2 LIMIT 1"
)

type DefaultMerchantPaymentGatewayConfigRepository struct {
	SqlConnection
}

func NewDefaultMerchantPaymentGatewayConfigRepository(sqlConnection SqlConnection) repository.MerchantPaymentGatewayConfigRepository {
	return &DefaultMerchantPaymentGatewayConfigRepository{SqlConnection: sqlConnection}
}

func (r *DefaultMerchantPaymentGatewayConfigRepository) FindByMerchantAndGateway(ctx context.Context, merchantId, gateway string) (*merchant.PaymentGatewayConfig, error) {
	var credential string
	found, err := r.Sql(sqlFindMerchantPaymentGatewayConfigByMerchantAndGateway).Values(merchantId, gateway).FindOne(ctx, &credential)
	if err != nil || !found {
		return nil, err
	}
	return &merchant.PaymentGatewayConfig{
		MerchantID:     merchantId,
		Credential:     credential,
		PaymentGateway: gateway,
	}, nil
}

func (r *DefaultMerchantPaymentGatewayConfigRepository) FindByMerchant(ctx context.Context, merchantId string) (*merchant.PaymentGatewayConfig, error) {
	var credential, gateway string
	found, err := r.Sql(sqlFindMerchantPaymentGatewayConfigByMerchant).Values(merchantId).FindOne(ctx, &credential, &gateway)
	if err != nil || !found {
		return nil, err
	}
	return &merchant.PaymentGatewayConfig{
		MerchantID:     merchantId,
		Credential:     credential,
		PaymentGateway: gateway,
	}, nil
}
