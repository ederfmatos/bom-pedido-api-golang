package token

type CustomerTokenManager interface {
	Encrypt(id string) (string, error)
}
