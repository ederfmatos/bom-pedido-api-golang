package factory

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/env"
	"bom-pedido-api/infra/token"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func tokenFactory(environment *env.Environment) *factory.TokenFactory {
	privateKey := loadPrivateKey(environment.JwePrivateKeyPath)
	tokenManager := token.NewCustomerTokenManager(privateKey)
	return factory.NewTokenFactory(tokenManager)
}

func loadPrivateKey(file string) *rsa.PrivateKey {
	pemData, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		panic("failed to decode PEM block containing private key")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	return key
}
