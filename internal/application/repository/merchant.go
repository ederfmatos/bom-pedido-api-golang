package repository

import (
	"bom-pedido-api/internal/domain/entity/merchant"
	"context"
)

type MerchantRepository interface {
	FindByTenantId(ctx context.Context, tenantId string) (*merchant.Merchant, error)
	IsActive(ctx context.Context, merchantId string) (bool, error)
	Create(ctx context.Context, merchant *merchant.Merchant) error
	Update(ctx context.Context, merchant *merchant.Merchant) error
}

type MerchantPaymentGatewayConfigRepository interface {
	FindByMerchantAndGateway(ctx context.Context, merchantId, gateway string) (*merchant.PaymentGatewayConfig, error)
	FindByMerchant(ctx context.Context, merchantId string) (*merchant.PaymentGatewayConfig, error)
	Create(ctx context.Context, config *merchant.PaymentGatewayConfig) error
}
