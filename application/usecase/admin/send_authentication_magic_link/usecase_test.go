package send_authentication_magic_link

import (
	"bom-pedido-api/domain/entity/admin"
	"bom-pedido-api/domain/entity/merchant"
	"bom-pedido-api/infra/event"
	"bom-pedido-api/infra/factory"
	"bom-pedido-api/infra/token"
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUseCase_Execute(t *testing.T) {
	applicationFactory := factory.NewTestApplicationFactory()
	baseUrl := faker.URL()
	tokenManager := token.NewFakeCustomerTokenManager()
	eventEmitter := event.NewMockEventHandler()
	applicationFactory.TokenManager = tokenManager
	applicationFactory.EventEmitter = eventEmitter

	t.Run("it should return nil if admin does not exists", func(t *testing.T) {
		useCase := New(baseUrl, applicationFactory)
		input := Input{Email: faker.Email()}
		err := useCase.Execute(context.Background(), input)
		require.Nil(t, err)
	})

	t.Run("it should return nil if merchant is inactive", func(t *testing.T) {
		ctx := context.Background()
		aMerchant, err := merchant.New(faker.Name(), faker.Email(), faker.Phonenumber(), faker.DomainName())
		require.Nil(t, err)
		aMerchant.Inactive()
		_ = applicationFactory.MerchantRepository.Create(ctx, aMerchant)

		anAdmin, _ := admin.New(faker.Name(), faker.Email(), aMerchant.Id)
		_ = applicationFactory.AdminRepository.Create(ctx, anAdmin)

		useCase := New(baseUrl, applicationFactory)
		input := Input{Email: anAdmin.GetEmail()}
		err = useCase.Execute(ctx, input)
		require.Nil(t, err)
	})

	t.Run("should return nil on success", func(t *testing.T) {
		ctx := context.Background()
		aMerchant, err := merchant.New(faker.Name(), faker.Email(), faker.Phonenumber(), faker.DomainName())
		require.Nil(t, err)
		_ = applicationFactory.MerchantRepository.Create(ctx, aMerchant)

		anAdmin, _ := admin.New(faker.Name(), faker.Email(), aMerchant.Id)
		_ = applicationFactory.AdminRepository.Create(context.Background(), anAdmin)

		magicLinkToken := faker.UUIDHyphenated()
		tokenManager.On("Encrypt", mock.Anything).Return(magicLinkToken, nil).Once()
		eventEmitter.On("Emit", mock.Anything, mock.Anything).Return(nil).Once()

		useCase := New(baseUrl, applicationFactory)
		input := Input{Email: anAdmin.GetEmail()}

		err = useCase.Execute(ctx, input)
		require.NoError(t, err)

		eventEmitter.AssertNumberOfCalls(t, "Emit", 1)
	})
}
