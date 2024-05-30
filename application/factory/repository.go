package factory

import (
	"bom-pedido-api/application/repository"
)

type RepositoryFactory struct {
	CustomerRepository repository.CustomerRepository
}

func NewRepositoryFactory(customerRepository repository.CustomerRepository) *RepositoryFactory {
	return &RepositoryFactory{CustomerRepository: customerRepository}
}
