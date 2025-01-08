package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/merchant"
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

			aMerchant, err := merchant.New(faker.Name(), faker.Email(), faker.Phonenumber(), faker.DomainName())
			require.NoError(t, err)

			savedMerchant, err := merchantRepository.FindByTenantId(ctx, aMerchant.TenantId)
			require.NoError(t, err)
			require.Nil(t, savedMerchant)

			merchantIsActive, err := merchantRepository.IsActive(ctx, aMerchant.Id)
			require.NoError(t, err)
			require.False(t, merchantIsActive)

			err = merchantRepository.Create(ctx, aMerchant)
			require.NoError(t, err)

			savedMerchant, err = merchantRepository.FindByTenantId(ctx, aMerchant.TenantId)
			require.NoError(t, err)
			require.NotNil(t, savedMerchant)
			require.Equal(t, aMerchant, savedMerchant)

			merchantIsActive, err = merchantRepository.IsActive(ctx, aMerchant.Id)
			require.NoError(t, err)
			require.True(t, merchantIsActive)

			aMerchant.Inactive()

			err = merchantRepository.Update(ctx, aMerchant)
			require.NoError(t, err)

			merchantIsActive, err = merchantRepository.IsActive(ctx, aMerchant.Id)
			require.NoError(t, err)
			require.False(t, merchantIsActive)
		})
	}
}
