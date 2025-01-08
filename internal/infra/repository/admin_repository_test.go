package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/admin"
	"bom-pedido-api/internal/domain/entity/merchant"
	"bom-pedido-api/internal/infra/test"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_AdminRepository(t *testing.T) {
	container := test.NewContainer()
	repositories := map[string]repository.AdminRepository{
		"AdminMemoryRepository": NewAdminMemoryRepository(),
		"AdminMongoRepository":  NewAdminMongoRepository(container.MongoDatabase()),
	}

	for name, adminRepository := range repositories {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			aMerchant, err := merchant.New(faker.Name(), faker.Email(), faker.Phonenumber(), faker.DomainName())
			require.NoError(t, err)

			aAdmin, err := admin.New(faker.Name(), faker.Email(), aMerchant.Id)
			require.NoError(t, err)

			savedAdmin, err := adminRepository.FindByEmail(ctx, aAdmin.GetEmail())
			require.NoError(t, err)
			require.Nil(t, savedAdmin)

			err = adminRepository.Create(ctx, aAdmin)
			require.NoError(t, err)

			savedAdmin, err = adminRepository.FindByEmail(ctx, aAdmin.GetEmail())
			require.NoError(t, err)
			require.NotNil(t, savedAdmin)
			require.Equal(t, aAdmin, savedAdmin)
		})
	}
}
