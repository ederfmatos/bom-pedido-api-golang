package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/infra/test"
	"bom-pedido-api/pkg/testify/require"
	"context"
	"testing"
)

func TestMerchantPaymentGatewayConfigRepository(t *testing.T) {
	container := test.NewContainer()

	repositories := map[string]repository.MerchantPaymentGatewayConfigRepository{
		"MerchantPaymentGatewayMongoRepository": NewMerchantPaymentGatewayConfigMongoRepository(container.MongoDatabase()),
	}

	for name, merchantPaymentGatewayConfigRepository := range repositories {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			merchantId := "123"
			paymentGatewayId := "456"
			merchantPaymentGatewayConfig, err := merchantPaymentGatewayConfigRepository.FindByMerchantAndGateway(ctx, merchantId, paymentGatewayId)
			require.NoError(t, err)
			require.Nil(t, merchantPaymentGatewayConfig)

			merchantPaymentGatewayConfig, err = merchantPaymentGatewayConfigRepository.FindByMerchant(ctx, merchantId)
			require.NoError(t, err)
			require.Nil(t, merchantPaymentGatewayConfig)

			merchantPaymentGatewayConfig = entity.NewMerchantPaymentGatewayConfig(merchantId, paymentGatewayId, "123")
			require.NoError(t, merchantPaymentGatewayConfigRepository.Create(ctx, merchantPaymentGatewayConfig))

			savedMerchantPaymentGatewayConfig, err := merchantPaymentGatewayConfigRepository.FindByMerchantAndGateway(ctx, merchantId, paymentGatewayId)
			require.NoError(t, err)
			require.NotNil(t, merchantPaymentGatewayConfig)
			require.Equal(t, merchantPaymentGatewayConfig, savedMerchantPaymentGatewayConfig)

			savedMerchantPaymentGatewayConfig, err = merchantPaymentGatewayConfigRepository.FindByMerchant(ctx, merchantId)
			require.NoError(t, err)
			require.NotNil(t, merchantPaymentGatewayConfig)
			require.Equal(t, merchantPaymentGatewayConfig, savedMerchantPaymentGatewayConfig)
		})
	}
}
