package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/customer"
	"bom-pedido-api/internal/domain/value_object"
	"bom-pedido-api/internal/infra/test"
	"context"
	"github.com/go-faker/faker/v4"
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

			aCustomerNotification := customer.NewNotification(value_object.NewID(), faker.Word())

			savedCustomerNotification, err := customerNotificationRepository.FindByCustomerId(ctx, aCustomerNotification.CustomerId)
			require.NoError(t, err)
			require.Nil(t, savedCustomerNotification)

			err = customerNotificationRepository.Upsert(ctx, aCustomerNotification)
			require.NoError(t, err)

			savedCustomerNotification, err = customerNotificationRepository.FindByCustomerId(ctx, aCustomerNotification.CustomerId)
			require.NoError(t, err)
			require.NotNil(t, savedCustomerNotification)
			require.Equal(t, aCustomerNotification, savedCustomerNotification)

			aCustomerNotification.Recipient = faker.Word()
			err = customerNotificationRepository.Upsert(ctx, aCustomerNotification)
			require.NoError(t, err)

			savedCustomerNotification, err = customerNotificationRepository.FindByCustomerId(ctx, aCustomerNotification.CustomerId)
			require.NoError(t, err)
			require.NotNil(t, savedCustomerNotification)
			require.Equal(t, aCustomerNotification, savedCustomerNotification)
		})
	}
}
