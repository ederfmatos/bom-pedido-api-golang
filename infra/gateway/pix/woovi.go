package pix

import (
	"bom-pedido-api/application/gateway"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/infra/http_client"
	"context"
	"time"
)

const (
	Woovi = "WOOVI"
)

type (
	wooviPixGateway struct {
		expirationInSeconds int
		httpClient          http_client.HTTPClient
	}

	wooviRefundOutput struct {
		Refunds []struct {
			Value  int64  `json:"value"`
			Status string `json:"status"`
		} `json:"refunds"`
	}

	wooviErrorResponse struct {
		ErrorMessage string `json:"error"`
	}

	wooviCreateChargeInput struct {
		CorrelationID string `json:"correlationID"`
		Value         int64  `json:"value"`
		Comment       string `json:"comment"`
		ExpiresIn     int    `json:"expiresIn"`
		Customer      struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"customer"`
	}

	wooviCreateChargeOutput struct {
		Charge *struct {
			CorrelationID  string    `json:"correlationID"`
			ExpiresDate    time.Time `json:"expiresDate"`
			BrCode         string    `json:"brCode"`
			PaymentLinkUrl string    `json:"paymentLinkUrl"`
		} `json:"charge"`
	}

	wooviGetChargeOutput struct {
		Charge *struct {
			Value          int64     `json:"value"`
			CorrelationID  string    `json:"correlationID"`
			ExpiresDate    time.Time `json:"expiresDate"`
			BrCode         string    `json:"brCode"`
			PaymentLinkUrl string    `json:"paymentLinkUrl"`
			Status         string    `json:"status"`
		} `json:"charge"`
	}

	wooviRefundChargeInput struct {
		CorrelationID string  `json:"correlationID"`
		Value         float64 `json:"value"`
	}

	wooviRefundChargeOutput struct {
		Refund struct {
			Status string `json:"status"`
		} `json:"refund"`
	}
)

func NewWooviPixGateway(
	httpClient http_client.HTTPClient,
	expirationTimeInMinutes int,
) gateway.PixGateway {
	return &wooviPixGateway{
		httpClient:          httpClient,
		expirationInSeconds: expirationTimeInMinutes * 60,
	}
}
func (e *wooviErrorResponse) Error() string {
	return e.ErrorMessage
}

func (g *wooviPixGateway) Name() string {
	return Woovi
}

func (g *wooviPixGateway) CreateQrCodePix(ctx context.Context, input gateway.CreateQrCodePixInput) (*gateway.CreateQrCodePixOutput, error) {
	paymentInput := wooviCreateChargeInput{
		CorrelationID: input.InternalOrderId,
		Value:         int64(input.Amount * 100),
		ExpiresIn:     g.expirationInSeconds,
		Comment:       input.Description,
		Customer: struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}{
			Name:  input.Customer.Name,
			Email: input.Customer.Email,
		},
	}
	response, err := g.httpClient.Post("/v1/charge").
		Body(paymentInput).
		Header("accept", "application/json").
		Header("content-type", "application/json").
		Header("Authorization", input.Credential).
		Execute(ctx)
	if err != nil {
		return nil, err
	}
	defer response.Close()
	if response.IsError() {
		return nil, response.ParseError(&wooviErrorResponse{})
	}
	var output wooviCreateChargeOutput
	if err = response.ParseBody(&output); err != nil {
		return nil, err
	}
	return &gateway.CreateQrCodePixOutput{
		PaymentGateway: Woovi,
		Id:             output.Charge.CorrelationID,
		QrCode:         output.Charge.BrCode,
		ExpiresAt:      output.Charge.ExpiresDate,
		QrCodeLink:     output.Charge.PaymentLinkUrl,
	}, nil
}

func (g *wooviPixGateway) GetPaymentById(ctx context.Context, input gateway.GetPaymentInput) (*gateway.GetPaymentOutput, error) {
	response, err := g.httpClient.Get("/v1/charge/", input.PaymentId).
		Header("accept", "application/json").
		Header("Authorization", input.Credential).
		Execute(ctx)
	if err != nil {
		return nil, err
	}
	defer response.Close()
	if response.IsError() {
		return nil, response.ParseError(&wooviErrorResponse{})
	}
	var output wooviGetChargeOutput
	if err = response.ParseBody(&output); err != nil {
		return nil, err
	}
	var status gateway.PaymentStatus
	switch output.Charge.Status {
	case "ACTIVE":
		status = gateway.TransactionPending
	case "EXPIRED":
		status = gateway.TransactionCancelled
	case "COMPLETED":
		status = gateway.TransactionPaid
		refundValue, err := g.getRefundValueFromCharge(ctx, input.PaymentId, input.Credential)
		if err != nil {
			return nil, err
		}
		if refundValue == output.Charge.Value {
			status = gateway.TransactionRefunded
		}
	default:
		return nil, nil
	}
	return &gateway.GetPaymentOutput{
		PaymentGateway: Woovi,
		Id:             output.Charge.CorrelationID,
		QrCode:         output.Charge.BrCode,
		ExpiresAt:      output.Charge.ExpiresDate,
		QrCodeLink:     output.Charge.PaymentLinkUrl,
		Status:         status,
	}, nil
}

func (g *wooviPixGateway) getRefundValueFromCharge(ctx context.Context, paymentId, credential string) (int64, error) {
	response, err := g.httpClient.Get("/v1/charge/", paymentId, "/refund").
		Header("accept", "application/json").
		Header("Authorization", credential).
		Execute(ctx)
	if err != nil {
		return 0, err
	}
	defer response.Close()
	if response.IsError() {
		return 0, response.ParseError(&wooviErrorResponse{})
	}
	var output wooviRefundOutput
	if err = response.ParseBody(&output); err != nil {
		return 0, err
	}
	value := int64(0)
	for _, refund := range output.Refunds {
		if refund.Status == "CONFIRMED" {
			value += refund.Value
		}
	}
	return value, nil
}

func (g *wooviPixGateway) RefundPix(ctx context.Context, input gateway.RefundPixInput) error {
	paymentInput := wooviRefundChargeInput{
		CorrelationID: input.PaymentId,
		Value:         input.Amount * 100,
	}
	response, err := g.httpClient.Post("/v1/charge/", input.PaymentId, "/refund").
		Header("accept", "application/json").
		Header("content-type", "application/json").
		Header("Authorization", input.Credential).
		Body(paymentInput).
		Execute(ctx)
	if err != nil {
		return err
	}
	defer response.Close()
	if response.IsError() {
		return response.ParseError(&wooviErrorResponse{})
	}
	var output wooviRefundChargeOutput
	if err = response.ParseBody(&output); err != nil {
		return err
	}
	switch output.Refund.Status {
	case "IN_PROCESSING", "CONFIRMED":
		return nil
	case "REJECTED":
		return errors.New("Refund rejected by payment gateway")
	default:
		return errors.New("Unexpected error")
	}
}
