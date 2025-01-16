package admin

import (
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/infra/event"
	"bom-pedido-api/internal/infra/factory"
	"bom-pedido-api/internal/infra/token"
	"bom-pedido-api/pkg/faker"
	"bom-pedido-api/pkg/testify/require"
	"context"
	"github.com/stretchr/testify/mock"
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
		useCase := NewSendAuthenticationMagicLink(baseUrl, applicationFactory)
		input := SendAuthenticationMagicLinkInput{Email: faker.Email()}
		err := useCase.Execute(context.Background(), input)
		require.Nil(t, err)
	})

	t.Run("it should return nil if merchant is inactive", func(t *testing.T) {
		ctx := context.Background()
		merchant, err := entity.NewMerchant(faker.Name(), faker.Email(), faker.PhoneNumber(), faker.DomainName())
		require.Nil(t, err)
		merchant.Inactive()
		_ = applicationFactory.MerchantRepository.Create(ctx, merchant)

		admin, _ := entity.NewAdmin(faker.Name(), faker.Email(), merchant.Id)
		_ = applicationFactory.AdminRepository.Create(ctx, admin)

		useCase := NewSendAuthenticationMagicLink(baseUrl, applicationFactory)
		input := SendAuthenticationMagicLinkInput{Email: admin.GetEmail()}
		err = useCase.Execute(ctx, input)
		require.Nil(t, err)
	})

	t.Run("should return nil on success", func(t *testing.T) {
		ctx := context.Background()
		merchant, err := entity.NewMerchant(faker.Name(), faker.Email(), faker.PhoneNumber(), faker.DomainName())
		require.Nil(t, err)
		_ = applicationFactory.MerchantRepository.Create(ctx, merchant)

		admin, _ := entity.NewAdmin(faker.Name(), faker.Email(), merchant.Id)
		_ = applicationFactory.AdminRepository.Create(context.Background(), admin)

		magicLinkToken := faker.Jwt()
		tokenManager.On("Encrypt", mock.Anything).Return(magicLinkToken, nil).Once()
		eventEmitter.On("Emit", mock.Anything, mock.Anything).Return(nil).Once()

		useCase := NewSendAuthenticationMagicLink(baseUrl, applicationFactory)
		input := SendAuthenticationMagicLinkInput{Email: admin.GetEmail()}

		err = useCase.Execute(ctx, input)
		require.NoError(t, err)

		eventEmitter.AssertNumberOfCalls(t, "Emit", 1)
	})
}
