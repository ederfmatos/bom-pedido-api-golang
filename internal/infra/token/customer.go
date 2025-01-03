package token

import (
	"bom-pedido-api/internal/application/token"
	"bom-pedido-api/internal/infra/json"
	"context"
	"crypto/rsa"
	"github.com/golang-jwt/jwe"
	"github.com/golang-jwt/jwt"
	"time"
)

type DefaultTokenManager struct {
	privateKey *rsa.PrivateKey
}

func NewCustomerTokenManager(privateKey *rsa.PrivateKey) token.Manager {
	return &DefaultTokenManager{privateKey: privateKey}
}

func (tokenManager *DefaultTokenManager) Encrypt(ctx context.Context, data token.Data) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Id:        data.Id,
		IssuedAt:  time.Now().Unix(),
		Issuer:    "BOM_PEDIDO_API",
		NotBefore: time.Now().Unix(),
		Subject:   data.TenantId,
		Audience:  data.Type,
	}
	tokenContent, err := json.Marshal(ctx, claims)
	if err != nil {
		return "", err
	}
	tokenJwe, err := jwe.NewJWE(jwe.KeyAlgorithmRSAOAEP, &tokenManager.privateKey.PublicKey, jwe.EncryptionTypeA256GCM, tokenContent)
	if err != nil {
		return "", err
	}
	return tokenJwe.CompactSerialize()
}

func (tokenManager *DefaultTokenManager) Decrypt(ctx context.Context, rawToken string) (*token.Data, error) {
	tokenContent, err := jwe.ParseEncrypted(rawToken)
	if err != nil {
		return nil, err
	}
	content, err := tokenContent.Decrypt(tokenManager.privateKey)
	if err != nil {
		return nil, err
	}
	var claims jwt.StandardClaims
	err = json.Unmarshal(ctx, content, &claims)
	if err != nil {
		return nil, err
	}
	err = claims.Valid()
	if err != nil {
		return nil, err
	}
	return &token.Data{
		Type:     claims.Audience,
		Id:       claims.Id,
		TenantId: claims.Subject,
	}, nil
}
