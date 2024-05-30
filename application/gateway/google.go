package gateway

type GoogleUser struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type GoogleGateway interface {
	GetUserByToken(token string) (*GoogleUser, error)
}
