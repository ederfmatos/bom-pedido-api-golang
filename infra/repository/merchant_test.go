package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/merchant"
	"bom-pedido-api/infra/test"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)

	savedMerchant, err := repository.FindByTenantId(ctx, aMerchant.TenantId)
	assert.NoError(t, err)
	assert.Nil(t, savedMerchant)

	merchantIsActive, err := repository.IsActive(ctx, aMerchant.Id)
	assert.NoError(t, err)
	assert.False(t, merchantIsActive)

	err = repository.Create(ctx, aMerchant)
	assert.NoError(t, err)

	savedMerchant, err = repository.FindByTenantId(ctx, aMerchant.TenantId)
	assert.NoError(t, err)
	assert.NotNil(t, savedMerchant)
	assert.Equal(t, aMerchant, savedMerchant)

	merchantIsActive, err = repository.IsActive(ctx, aMerchant.Id)
	assert.NoError(t, err)
	assert.True(t, merchantIsActive)

	aMerchant.Inactive()

	err = repository.Update(ctx, aMerchant)
	assert.NoError(t, err)

	merchantIsActive, err = repository.IsActive(ctx, aMerchant.Id)
	assert.NoError(t, err)
	assert.False(t, merchantIsActive)
}
