package pix

import (
	"bom-pedido-api/application/gateway"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/infra/json"
	"bytes"
	"context"
	"net/http"
	"time"
)

const (
	woovi = "WOOVI"
)

type (
	WooviPixGateway struct {
		expirationInSeconds                    int
		merchantPaymentGatewayConfigRepository repository.MerchantPaymentGatewayConfigRepository
		url                                    string
	}

	WooviErrorResponse struct {
		Error string `json:"error"`
	}

	WooviCreateChargeInput struct {
		CorrelationID string  `json:"correlationID"`
		Value         float64 `json:"value"`
		Comment       string  `json:"comment"`
		ExpiresIn     int     `json:"expiresIn"`
		Customer      struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"customer"`
	}

	WooviCreateChargeOutput struct {
		Charge *struct {
			CorrelationID  string    `json:"correlationID"`
			ExpiresDate    time.Time `json:"expiresDate"`
			BrCode         string    `json:"brCode"`
			PaymentLinkUrl string    `json:"paymentLinkUrl"`
		} `json:"charge"`
	}

	WooviGetChargeOutput struct {
		Charge *struct {
			CorrelationID  string    `json:"correlationID"`
			ExpiresDate    time.Time `json:"expiresDate"`
			BrCode         string    `json:"brCode"`
			PaymentLinkUrl string    `json:"paymentLinkUrl"`
			Status         string    `json:"status"`
		} `json:"charge"`
	}

	WooviRefundChargeInput struct {
		CorrelationID string  `json:"correlationID"`
		Value         float64 `json:"value"`
	}

	WooviRefundChargeOutput struct {
		Refund struct {
			Status string `json:"status"`
		} `json:"refund"`
	}
)

func NewWooviPixGateway(expirationTimeInMinutes int, merchantPaymentGatewayConfigRepository repository.MerchantPaymentGatewayConfigRepository) gateway.PixGateway {
	return &WooviPixGateway{
		url:                                    "https://api.openpix.com.br/api",
		expirationInSeconds:                    expirationTimeInMinutes * 60,
		merchantPaymentGatewayConfigRepository: merchantPaymentGatewayConfigRepository,
	}
}

func (g *WooviPixGateway) Name() string {
	return woovi
}

func (g *WooviPixGateway) getMerchantAccessToken(ctx context.Context, merchantId string) (*string, error) {
	gatewayConfig, err := g.merchantPaymentGatewayConfigRepository.FindByMerchantAndGateway(ctx, merchantId, woovi)
	if err != nil {
		return nil, err
	}
	if gatewayConfig == nil {
		return nil, errors.New("gateway config not found")
	}
	return &gatewayConfig.AccessToken, nil
}

func (g *WooviPixGateway) CreateQrCodePix(ctx context.Context, input gateway.CreateQrCodePixInput) (*gateway.CreateQrCodePixOutput, error) {
	accessToken, err := g.getMerchantAccessToken(ctx, input.MerchantId)
	if err != nil {
		return nil, err
	}
	paymentInput := WooviCreateChargeInput{
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
	paymentBytes, err := json.Marshal(ctx, paymentInput)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", g.url+"/v1/charge", bytes.NewBuffer(paymentBytes))
	if err != nil {
		return nil, err
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("authorization", *accessToken)
	request.Header.Add("content-type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		errorResponse := WooviErrorResponse{Error: response.Status}
		if err = json.Decode(ctx, response.Body, &errorResponse); err != nil {
			return nil, err
		}
		return nil, errors.New(errorResponse.Error)
	}
	var createPaymentOutput WooviCreateChargeOutput
	if err = json.Decode(ctx, response.Body, &createPaymentOutput); err != nil {
		return nil, err
	}
	return &gateway.CreateQrCodePixOutput{
		PaymentGateway: woovi,
		Id:             createPaymentOutput.Charge.CorrelationID,
		QrCode:         createPaymentOutput.Charge.BrCode,
		ExpiresAt:      createPaymentOutput.Charge.ExpiresDate,
		QrCodeLink:     createPaymentOutput.Charge.PaymentLinkUrl,
	}, nil
}

func (g *WooviPixGateway) GetPaymentById(ctx context.Context, merchantId, paymentId string) (*gateway.GetPaymentOutput, error) {
	accessToken, err := g.getMerchantAccessToken(ctx, merchantId)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, "GET", g.url+"/v1/charge/"+paymentId, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("authorization", *accessToken)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		errorResponse := WooviErrorResponse{Error: response.Status}
		if err = json.Decode(ctx, response.Body, &errorResponse); err != nil {
			return nil, err
		}
		return nil, errors.New(errorResponse.Error)
	}
	var output WooviGetChargeOutput
	if err = json.Decode(ctx, response.Body, &output); err != nil {
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

func (g *WooviPixGateway) RefundPix(ctx context.Context, input gateway.RefundPixInput) error {
	accessToken, err := g.getMerchantAccessToken(ctx, input.MerchantId)
	if err != nil {
		return err
	}
	paymentInput := WooviRefundChargeInput{
		CorrelationID: input.PaymentId,
		Value:         input.Amount * 100,
	}
	paymentBytes, err := json.Marshal(ctx, paymentInput)
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", g.url+"/v1/charge", bytes.NewBuffer(paymentBytes))
	if err != nil {
		return err
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("authorization", *accessToken)
	request.Header.Add("content-type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		errorResponse := WooviErrorResponse{Error: response.Status}
		if err = json.Decode(ctx, response.Body, &errorResponse); err != nil {
			return err
		}
		return errors.New(errorResponse.Error)
	}
	var output WooviRefundChargeOutput
	if err = json.Decode(ctx, response.Body, &output); err != nil {
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
