package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"context"
)

type MerchantRepository interface {
	FindByTenantId(ctx context.Context, tenantId string) (*entity.Merchant, error)
	IsActive(ctx context.Context, merchantId string) (bool, error)
	Create(ctx context.Context, merchant *entity.Merchant) error
	Update(ctx context.Context, merchant *entity.Merchant) error
}

type MerchantPaymentGatewayConfigRepository interface {
	FindByMerchantAndGateway(ctx context.Context, merchantId, gateway string) (*entity.MerchantPaymentGatewayConfig, error)
	FindByMerchant(ctx context.Context, merchantId string) (*entity.MerchantPaymentGatewayConfig, error)
	Create(ctx context.Context, config *entity.MerchantPaymentGatewayConfig) error
}
