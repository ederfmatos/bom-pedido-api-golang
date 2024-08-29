package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/admin"
	"bom-pedido-api/domain/entity/merchant"
	"bom-pedido-api/infra/test"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_AdminSqlRepository(t *testing.T) {
	container := test.NewContainer()
	sqlConnection := NewDefaultSqlConnection(container.Database)
	adminSqlRepository := NewDefaultAdminRepository(sqlConnection)
	merchantRepository := NewDefaultMerchantRepository(sqlConnection)
	runAdminTests(t, adminSqlRepository, merchantRepository)
}

func Test_AdminMemoryRepository(t *testing.T) {
	adminRepository := NewAdminMemoryRepository()
	merchantRepository := NewMerchantMemoryRepository()
	runAdminTests(t, adminRepository, merchantRepository)
}

func runAdminTests(t *testing.T, repository repository.AdminRepository, merchantRepository repository.MerchantRepository) {
	ctx := context.TODO()

	aMerchant, err := merchant.New(faker.Name(), faker.Email(), faker.Phonenumber(), faker.DomainName())
	require.NoError(t, err)

	err = merchantRepository.Create(ctx, aMerchant)
	require.NoError(t, err)

	aAdmin, err := admin.New(faker.Name(), faker.Email(), aMerchant.Id)
	require.NoError(t, err)

	savedAdmin, err := repository.FindByEmail(ctx, aAdmin.GetEmail())
	require.NoError(t, err)
	require.Nil(t, savedAdmin)

	err = repository.Create(ctx, aAdmin)
	require.NoError(t, err)

	savedAdmin, err = repository.FindByEmail(ctx, aAdmin.GetEmail())
	require.NoError(t, err)
	require.NotNil(t, savedAdmin)
	require.Equal(t, aAdmin, savedAdmin)
}
