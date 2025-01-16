package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/infra/test"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_MerchantRepository(t *testing.T) {
	container := test.NewContainer()
	repositories := map[string]repository.MerchantRepository{
		"MerchantMemoryRepository": NewMerchantMemoryRepository(),
		"MerchantMongoRepository":  NewMerchantMongoRepository(container.MongoDatabase()),
	}

	for name, merchantRepository := range repositories {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			merchant, err := entity.NewMerchant(faker.Name(), faker.Email(), faker.Phonenumber(), faker.DomainName())
			require.NoError(t, err)

			savedMerchant, err := merchantRepository.FindByTenantId(ctx, merchant.TenantId)
			require.NoError(t, err)
			require.Nil(t, savedMerchant)

			merchantIsActive, err := merchantRepository.IsActive(ctx, merchant.Id)
			require.NoError(t, err)
			require.False(t, merchantIsActive)

			err = merchantRepository.Create(ctx, merchant)
			require.NoError(t, err)

			savedMerchant, err = merchantRepository.FindByTenantId(ctx, merchant.TenantId)
			require.NoError(t, err)
			require.NotNil(t, savedMerchant)
			require.Equal(t, merchant, savedMerchant)

			merchantIsActive, err = merchantRepository.IsActive(ctx, merchant.Id)
			require.NoError(t, err)
			require.True(t, merchantIsActive)

			merchant.Inactive()

			err = merchantRepository.Update(ctx, merchant)
			require.NoError(t, err)

			merchantIsActive, err = merchantRepository.IsActive(ctx, merchant.Id)
			require.NoError(t, err)
			require.False(t, merchantIsActive)
		})
	}
}
