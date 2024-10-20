package factory

import (
	"bom-pedido-api/internal/application/token"
)

type TokenFactory struct {
	TokenManager token.Manager
}

func NewTokenFactory(tokenManager token.Manager) *TokenFactory {
	return &TokenFactory{TokenManager: tokenManager}
}
