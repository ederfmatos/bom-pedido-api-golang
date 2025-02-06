package token

import (
	"bom-pedido-api/internal/application/token"
	"context"
	"crypto/rsa"
	"encoding/json"
	"github.com/golang-jwt/jwe"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Manager struct {
	privateKey *rsa.PrivateKey
}

func NewTokenManager(privateKey *rsa.PrivateKey) *Manager {
	return &Manager{privateKey: privateKey}
}

func (m *Manager) Encrypt(ctx context.Context, data token.Data) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		ID:        data.Id,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "BOM_PEDIDO_API",
		NotBefore: jwt.NewNumericDate(time.Now()),
		Subject:   data.TenantId,
		Audience:  []string{data.Type},
	}
	tokenContent, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	tokenJwe, err := jwe.NewJWE(jwe.KeyAlgorithmRSAOAEP, &m.privateKey.PublicKey, jwe.EncryptionTypeA256GCM, tokenContent)
	if err != nil {
		return "", err
	}
	return tokenJwe.CompactSerialize()
}

func (m *Manager) Decrypt(ctx context.Context, rawToken string) (*token.Data, error) {
	tokenContent, err := jwe.ParseEncrypted(rawToken)
	if err != nil {
		return nil, err
	}
	content, err := tokenContent.Decrypt(m.privateKey)
	if err != nil {
		return nil, err
	}
	var claims jwt.RegisteredClaims
	err = json.Unmarshal(content, &claims)
	if err != nil {
		return nil, err
	}
	//err = claims.Valid()
	//if err != nil {
	//	return nil, err
	//}
	return &token.Data{
		Type:     claims.Audience[0],
		Id:       claims.ID,
		TenantId: claims.Subject,
	}, nil
}
