package gateway

type EmailGateway interface {
	Send(to, subject, template string, data map[string]string) error
}
