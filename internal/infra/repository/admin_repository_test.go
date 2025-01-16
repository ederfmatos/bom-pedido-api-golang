package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/infra/test"
	"bom-pedido-api/pkg/faker"
	"context"
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

			merchant, err := entity.NewMerchant(faker.Name(), faker.Email(), faker.PhoneNumber(), faker.DomainName())
			require.NoError(t, err)

			aAdmin, err := entity.NewAdmin(faker.Name(), faker.Email(), merchant.Id)
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
