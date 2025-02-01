package factory

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/infra/config"
	"bom-pedido-api/internal/infra/token"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func tokenFactory(environment *config.Environment) (*factory.TokenFactory, error) {
	privateKey, err := loadPrivateKey(environment.JwePrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("load private key: %v", err)
	}

	tokenManager := token.NewTokenManager(privateKey)
	return factory.NewTokenFactory(tokenManager), nil
}

func loadPrivateKey(file string) (*rsa.PrivateKey, error) {
	pemData, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read pem file: %v", err)
	}

	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %v", err)
	}

	return key, nil
}
