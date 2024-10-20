package token

import "context"

type (
	Data struct {
		Type     string
		Id       string
		TenantId string
	}

	Manager interface {
		Encrypt(ctx context.Context, data Data) (string, error)
		Decrypt(ctx context.Context, token string) (*Data, error)
	}
)
