package factory

import (
	"bom-pedido-api/application/token"
)

type TokenFactory struct {
	TokenManager token.Manager
}

func NewTokenFactory(tokenManager token.Manager) *TokenFactory {
	return &TokenFactory{TokenManager: tokenManager}
}
