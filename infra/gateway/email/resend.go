package email

import (
	"bom-pedido-api/application/gateway"
	"github.com/resend/resend-go/v2"
)

type ResendEmailGateway struct {
	client         *resend.Client
	templateLoader TemplateLoader
	from           string
}

func NewResendEmailGateway(templateLoader TemplateLoader, from, apiKey string) gateway.EmailGateway {
	client := resend.NewClient(apiKey)
	return &ResendEmailGateway{client: client, templateLoader: templateLoader, from: from}
}

func (r *ResendEmailGateway) Send(to, subject, template string, data map[string]string) error {
	html, err := r.templateLoader.Load(template, data)
	if err != nil {
		return err
	}
	params := &resend.SendEmailRequest{
		To:      []string{to},
		From:    r.from,
		Html:    html,
		Subject: subject,
	}
	_, err = r.client.Emails.Send(params)
	return err
}
