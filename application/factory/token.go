package factory

import (
	"bom-pedido-api/application/token"
)

type TokenFactory struct {
	CustomerTokenManager token.CustomerTokenManager
}

func NewTokenFactory(customerTokenManager token.CustomerTokenManager) *TokenFactory {
	return &TokenFactory{CustomerTokenManager: customerTokenManager}
}
