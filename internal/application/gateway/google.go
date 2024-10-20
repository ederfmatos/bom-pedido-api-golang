package gateway

import "context"

type GoogleUser struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type GoogleGateway interface {
	GetUserByToken(ctx context.Context, token string) (*GoogleUser, error)
}
