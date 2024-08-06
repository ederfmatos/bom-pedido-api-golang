package token

import "context"

type CustomerTokenManager interface {
	Encrypt(ctx context.Context, id string) (string, error)
	Decrypt(ctx context.Context, token string) (string, error)
}
