package token

import (
	"bom-pedido-api/infra/json"
	"context"
	"crypto/rsa"
	"github.com/golang-jwt/jwe"
	"github.com/golang-jwt/jwt"
	"time"
)

type CustomerTokenManager struct {
	privateKey *rsa.PrivateKey
}

func NewCustomerTokenManager(privateKey *rsa.PrivateKey) *CustomerTokenManager {
	return &CustomerTokenManager{privateKey: privateKey}
}

func (tokenManager *CustomerTokenManager) Encrypt(ctx context.Context, id string) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Id:        id,
		IssuedAt:  time.Now().Unix(),
		Issuer:    "BOM_PEDIDO_API",
		NotBefore: time.Now().Unix(),
		Subject:   id,
	}
	tokenContent, err := json.Marshal(ctx, claims)
	if err != nil {
		return "", err
	}
	token, err := jwe.NewJWE(jwe.KeyAlgorithmRSAOAEP, &tokenManager.privateKey.PublicKey, jwe.EncryptionTypeA256GCM, tokenContent)
	if err != nil {
		return "", err
	}
	return token.CompactSerialize()
}

func (tokenManager *CustomerTokenManager) Decrypt(ctx context.Context, token string) (string, error) {
	tokenContent, err := jwe.ParseEncrypted(token)
	if err != nil {
		return "", err
	}
	content, err := tokenContent.Decrypt(tokenManager.privateKey)
	if err != nil {
		return "", err
	}
	var claims jwt.StandardClaims
	err = json.Unmarshal(ctx, content, &claims)
	if err != nil {
		return "", err
	}
	err = claims.Valid()
	if err != nil {
		return "", err
	}
	return claims.Id, nil
}
