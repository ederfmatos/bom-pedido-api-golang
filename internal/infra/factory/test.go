package factory

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/infra/event"
	"bom-pedido-api/internal/infra/gateway/email"
	"bom-pedido-api/internal/infra/gateway/google"
	"bom-pedido-api/internal/infra/gateway/notification"
	"bom-pedido-api/internal/infra/gateway/pix"
	"bom-pedido-api/internal/infra/lock"
	"bom-pedido-api/internal/infra/repository"
	"bom-pedido-api/internal/infra/token"
)

func NewTestApplicationFactory() *factory.ApplicationFactory {
	return factory.NewApplicationFactory(
		factory.NewGatewayFactory(
			google.NewFakeGoogleGateway(),
			pix.NewFakePixGateway(),
			notification.NewMockNotificationGateway(),
			email.NewFakeEmailGateway(),
		),
		factory.NewRepositoryFactory(
			repository.NewCustomerMemoryRepository(),
			repository.NewProductMemoryRepository(),
			repository.NewShoppingCartMemoryRepository(),
			repository.NewOrderMemoryRepository(),
			repository.NewAdminMemoryRepository(),
			repository.NewMerchantMemoryRepository(),
			repository.NewTransactionMemoryRepository(),
			repository.NewOrderStatusHistoryMemoryRepository(),
			repository.NewCustomerNotificationMemoryRepository(),
			repository.NewNotificationMemoryRepository(),
			repository.NewCategoriesMemoryRepository(),
		),
		factory.NewTokenFactory(token.NewFakeCustomerTokenManager()),
		factory.NewEventFactory(event.NewMemoryEventHandler()),
		nil,
		lock.NewMemoryLocker(),
	)
}
