package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/merchant"
	"bom-pedido-api/infra/test"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_MerchantSqlRepository(t *testing.T) {
	container := test.NewContainer()
	sqlConnection := NewDefaultSqlConnection(container.Database)
	merchantSqlRepository := NewDefaultMerchantRepository(sqlConnection)
	runMerchantTests(t, merchantSqlRepository)
}

func Test_MerchantMemoryRepository(t *testing.T) {
	merchantSqlRepository := NewMerchantMemoryRepository()
	runMerchantTests(t, merchantSqlRepository)
}

func runMerchantTests(t *testing.T, repository repository.MerchantRepository) {
	ctx := context.TODO()

	aMerchant, err := merchant.New(faker.Name(), faker.Email(), faker.Phonenumber(), faker.DomainName())
	require.NoError(t, err)

	savedMerchant, err := repository.FindByTenantId(ctx, aMerchant.TenantId)
	require.NoError(t, err)
	require.Nil(t, savedMerchant)

	merchantIsActive, err := repository.IsActive(ctx, aMerchant.Id)
	require.NoError(t, err)
	require.False(t, merchantIsActive)

	err = repository.Create(ctx, aMerchant)
	require.NoError(t, err)

	savedMerchant, err = repository.FindByTenantId(ctx, aMerchant.TenantId)
	require.NoError(t, err)
	require.NotNil(t, savedMerchant)
	require.Equal(t, aMerchant, savedMerchant)

	merchantIsActive, err = repository.IsActive(ctx, aMerchant.Id)
	require.NoError(t, err)
	require.True(t, merchantIsActive)

	aMerchant.Inactive()

	err = repository.Update(ctx, aMerchant)
	require.NoError(t, err)

	merchantIsActive, err = repository.IsActive(ctx, aMerchant.Id)
	require.NoError(t, err)
	require.False(t, merchantIsActive)
}
