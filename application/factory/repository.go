package factory

import (
	"bom-pedido-api/application/repository"
)

type RepositoryFactory struct {
	CustomerRepository     repository.CustomerRepository
	ProductRepository      repository.ProductRepository
	OrderRepository        repository.OrderRepository
	ShoppingCartRepository repository.ShoppingCartRepository
	AdminRepository        repository.AdminRepository
}

func NewRepositoryFactory(
	customerRepository repository.CustomerRepository,
	productRepository repository.ProductRepository,
	shoppingCartRepository repository.ShoppingCartRepository,
	orderRepository repository.OrderRepository,
	adminRepository repository.AdminRepository,
) *RepositoryFactory {
	return &RepositoryFactory{
		CustomerRepository:     customerRepository,
		ProductRepository:      productRepository,
		ShoppingCartRepository: shoppingCartRepository,
		OrderRepository:        orderRepository,
		AdminRepository:        adminRepository,
	}
}
