package pix

import (
	"bom-pedido-api/application/gateway"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/infra/http_client"
	"context"
	"time"
)

const (
	woovi = "WOOVI"
)

type (
	wooviPixGateway struct {
		expirationInSeconds                    int
		merchantPaymentGatewayConfigRepository repository.MerchantPaymentGatewayConfigRepository
		httpClient                             http_client.HttpClient
	}

	wooviErrorResponse struct {
		Error string `json:"error"`
	}

	wooviCreateChargeInput struct {
		CorrelationID string  `json:"correlationID"`
		Value         float64 `json:"value"`
		Comment       string  `json:"comment"`
		ExpiresIn     int     `json:"expiresIn"`
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
	expirationTimeInMinutes int,
	merchantPaymentGatewayConfigRepository repository.MerchantPaymentGatewayConfigRepository,
	httpClient http_client.HttpClient,
) gateway.PixGateway {
	return &wooviPixGateway{
		httpClient:                             httpClient,
		expirationInSeconds:                    expirationTimeInMinutes * 60,
		merchantPaymentGatewayConfigRepository: merchantPaymentGatewayConfigRepository,
	}
}

func (g *wooviPixGateway) Name() string {
	return woovi
}

func (g *wooviPixGateway) getMerchantAccessToken(ctx context.Context, merchantId string) (*string, error) {
	gatewayConfig, err := g.merchantPaymentGatewayConfigRepository.FindByMerchantAndGateway(ctx, merchantId, woovi)
	if err != nil {
		return nil, err
	}
	if gatewayConfig == nil {
		return nil, errors.New("gateway config not found")
	}
	return &gatewayConfig.AccessToken, nil
}

func (g *wooviPixGateway) CreateQrCodePix(ctx context.Context, input gateway.CreateQrCodePixInput) (*gateway.CreateQrCodePixOutput, error) {
	accessToken, err := g.getMerchantAccessToken(ctx, input.MerchantId)
	if err != nil {
		return nil, err
	}
	paymentInput := wooviCreateChargeInput{
		CorrelationID: input.InternalOrderId,
		Value:         input.Amount * 100,
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
		Header("Authorization", *accessToken).
		Execute(ctx)
	if err != nil {
		return nil, err
	}
	if response.IsError() {
		var wooviError wooviErrorResponse
		if err = response.ParseBody(&wooviError); err != nil {
			return nil, err
		}
		return nil, errors.New(wooviError.Error)
	}
	var output wooviCreateChargeOutput
	if err = response.ParseBody(&output); err != nil {
		return nil, err
	}
	return &gateway.CreateQrCodePixOutput{
		PaymentGateway: woovi,
		Id:             output.Charge.CorrelationID,
		QrCode:         output.Charge.BrCode,
		ExpiresAt:      output.Charge.ExpiresDate,
		QrCodeLink:     output.Charge.PaymentLinkUrl,
	}, nil
}

func (g *wooviPixGateway) GetPaymentById(ctx context.Context, merchantId, paymentId string) (*gateway.GetPaymentOutput, error) {
	accessToken, err := g.getMerchantAccessToken(ctx, merchantId)
	if err != nil {
		return nil, err
	}
	response, err := g.httpClient.Get("/v1/charge/", paymentId).
		Header("accept", "application/json").
		Header("Authorization", *accessToken).
		Execute(ctx)
	if err != nil {
		return nil, err
	}
	if response.IsError() {
		var wooviError wooviErrorResponse
		if err = response.ParseBody(&wooviError); err != nil {
			return nil, err
		}
		return nil, errors.New(wooviError.Error)
	}
	var output wooviGetChargeOutput
	if err = response.ParseBody(&output); err != nil {
		return nil, err
	}
	var status gateway.PaymentStatus
	switch output.Charge.Status {
	case "ACTIVE":
		status = gateway.TransactionPending
		break
	case "EXPIRED":
		status = gateway.TransactionCancelled
		break
	case "COMPLETED":
		status = gateway.TransactionPaid
		break
	default:
		return nil, nil
	}
	return &gateway.GetPaymentOutput{
		PaymentGateway: woovi,
		Id:             output.Charge.CorrelationID,
		QrCode:         output.Charge.BrCode,
		ExpiresAt:      output.Charge.ExpiresDate,
		QrCodeLink:     output.Charge.PaymentLinkUrl,
		Status:         status,
	}, nil
}

func (g *wooviPixGateway) RefundPix(ctx context.Context, input gateway.RefundPixInput) error {
	accessToken, err := g.getMerchantAccessToken(ctx, input.MerchantId)
	if err != nil {
		return err
	}
	paymentInput := wooviRefundChargeInput{
		CorrelationID: input.PaymentId,
		Value:         input.Amount * 100,
	}
	response, err := g.httpClient.Post("/v1/charge/", input.PaymentId, "/refund").
		Header("accept", "application/json").
		Header("content-type", "application/json").
		Header("Authorization", *accessToken).
		Body(paymentInput).
		Execute(ctx)
	if err != nil {
		return err
	}
	if response.IsError() {
		var wooviError wooviErrorResponse
		if err = response.ParseBody(&wooviError); err != nil {
			return err
		}
		return errors.New(wooviError.Error)
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
