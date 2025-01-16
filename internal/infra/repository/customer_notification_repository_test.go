package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/test"
	"bom-pedido-api/pkg/faker"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCustomerNotificationRepository(t *testing.T) {
	container := test.NewContainer()

	repositories := map[string]repository.CustomerNotificationRepository{
		"CustomerNotificationMongoRepository": NewCustomerNotificationMongoRepository(container.MongoDatabase()),
	}

	for name, customerNotificationRepository := range repositories {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			customerNotification := entity.NewCustomerNotification(value_object.NewID(), faker.Word())

			savedCustomerNotification, err := customerNotificationRepository.FindByCustomerId(ctx, customerNotification.CustomerId)
			require.NoError(t, err)
			require.Nil(t, savedCustomerNotification)

			err = customerNotificationRepository.Upsert(ctx, customerNotification)
			require.NoError(t, err)

			savedCustomerNotification, err = customerNotificationRepository.FindByCustomerId(ctx, customerNotification.CustomerId)
			require.NoError(t, err)
			require.NotNil(t, savedCustomerNotification)
			require.Equal(t, customerNotification, savedCustomerNotification)

			customerNotification.Recipient = faker.Word()
			err = customerNotificationRepository.Upsert(ctx, customerNotification)
			require.NoError(t, err)

			savedCustomerNotification, err = customerNotificationRepository.FindByCustomerId(ctx, customerNotification.CustomerId)
			require.NoError(t, err)
			require.NotNil(t, savedCustomerNotification)
			require.Equal(t, customerNotification, savedCustomerNotification)
		})
	}
}
