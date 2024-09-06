package factory

import (
	"bom-pedido-api/application/repository"
)

type RepositoryFactory struct {
	CustomerRepository             repository.CustomerRepository
	ProductRepository              repository.ProductRepository
	OrderRepository                repository.OrderRepository
	ShoppingCartRepository         repository.ShoppingCartRepository
	AdminRepository                repository.AdminRepository
	MerchantRepository             repository.MerchantRepository
	TransactionRepository          repository.TransactionRepository
	OrderHistoryRepository         repository.OrderStatusHistoryRepository
	CustomerNotificationRepository repository.CustomerNotificationRepository
	NotificationRepository         repository.NotificationRepository
}

func NewRepositoryFactory(
	customerRepository repository.CustomerRepository,
	productRepository repository.ProductRepository,
	shoppingCartRepository repository.ShoppingCartRepository,
	orderRepository repository.OrderRepository,
	adminRepository repository.AdminRepository,
	merchantRepository repository.MerchantRepository,
	transactionRepository repository.TransactionRepository,
	orderHistoryRepository repository.OrderStatusHistoryRepository,
	customerNotificationRepository repository.CustomerNotificationRepository,
	notificationRepository repository.NotificationRepository,
) *RepositoryFactory {
	return &RepositoryFactory{
		CustomerRepository:             customerRepository,
		ProductRepository:              productRepository,
		ShoppingCartRepository:         shoppingCartRepository,
		OrderRepository:                orderRepository,
		AdminRepository:                adminRepository,
		MerchantRepository:             merchantRepository,
		TransactionRepository:          transactionRepository,
		OrderHistoryRepository:         orderHistoryRepository,
		CustomerNotificationRepository: customerNotificationRepository,
		NotificationRepository:         notificationRepository,
	}
}
