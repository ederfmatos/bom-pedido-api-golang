package pix

import (
	"bom-pedido-api/application/gateway"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/errors"
	"context"
	"fmt"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
	"github.com/mercadopago/sdk-go/pkg/refund"
	"strconv"
	"time"
)

const (
	mercadoPago     = "MERCADO_PAGO"
	applicationName = "bom-pedido"
)

type MercadoPagoPixGateway struct {
	notificationUrl                        string
	expirationTimeInMinutes                int
	merchantPaymentGatewayConfigRepository repository.MerchantPaymentGatewayConfigRepository
}

func NewMercadoPagoPixGateway(
	notificationUrl string,
	expirationTimeInMinutes int,
	merchantPaymentGatewayConfigRepository repository.MerchantPaymentGatewayConfigRepository,
) gateway.PixGateway {
	return &MercadoPagoPixGateway{
		notificationUrl:                        notificationUrl,
		expirationTimeInMinutes:                expirationTimeInMinutes,
		merchantPaymentGatewayConfigRepository: merchantPaymentGatewayConfigRepository,
	}
}

func (g *MercadoPagoPixGateway) Name() string {
	return mercadoPago
}

func (g *MercadoPagoPixGateway) getMerchantConfig(ctx context.Context, merchantId string) (*string, error) {
	gatewayConfig, err := g.merchantPaymentGatewayConfigRepository.FindByMerchantAndGateway(ctx, merchantId, mercadoPago)
	if err != nil {
		return nil, err
	}
	if gatewayConfig == nil {
		return nil, errors.New("gateway config not found")
	}
	return &gatewayConfig.AccessToken, nil
}

func (g *MercadoPagoPixGateway) getConfig(ctx context.Context, merchantId string) (*config.Config, error) {
	gatewayConfig, err := g.merchantPaymentGatewayConfigRepository.FindByMerchantAndGateway(ctx, merchantId, mercadoPago)
	if err != nil {
		return nil, err
	}
	if gatewayConfig == nil {
		return nil, gateway.MerchantGatewayConfigNotFoundError
	}
	if err != nil {
		return nil, err
	}
	cfg, err := config.New(gatewayConfig.AccessToken, config.WithPlatformID(applicationName))
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (g *MercadoPagoPixGateway) CreateQrCodePix(ctx context.Context, input gateway.CreateQrCodePixInput) (*gateway.CreateQrCodePixOutput, error) {
	cfg, err := g.getConfig(ctx, input.MerchantId)
	if err != nil {
		return nil, err
	}
	client := payment.NewClient(cfg)
	expiresAt := time.Now().Add(time.Minute * time.Duration(g.expirationTimeInMinutes)).Truncate(time.Millisecond)
	request := payment.Request{
		TransactionAmount: input.Amount,
		Description:       input.Description,
		PaymentMethodID:   "pix",
		DateOfExpiration:  &expiresAt,
		NotificationURL:   fmt.Sprintf("%s/%s/%s", g.notificationUrl, mercadoPago, input.InternalOrderId),
		Payer: &payment.PayerRequest{
			FirstName: input.Customer.Name,
			Email:     input.Customer.Email,
		},
		Metadata: map[string]any{
			"orderId": input.InternalOrderId,
		},
	}
	resource, err := client.Create(ctx, request)
	if err != nil {
		return nil, err
	}
	return &gateway.CreateQrCodePixOutput{
		Id:             strconv.Itoa(resource.ID),
		QrCode:         resource.PointOfInteraction.TransactionData.QRCode,
		ExpiresAt:      expiresAt,
		PaymentGateway: mercadoPago,
		QrCodeLink:     resource.PointOfInteraction.TransactionData.TicketURL,
	}, nil
}

func (g *MercadoPagoPixGateway) GetPaymentById(ctx context.Context, merchantId, id string) (*gateway.GetPaymentOutput, error) {
	cfg, err := g.getConfig(ctx, merchantId)
	if err != nil {
		return nil, err
	}
	client := payment.NewClient(cfg)
	paymentId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	paymentResponse, err := client.Get(ctx, paymentId)
	if err != nil {
		return nil, err
	}
	var status gateway.PaymentStatus
	switch paymentResponse.Status {
	case "pending", "authorized", "in_process":
		status = gateway.TransactionPending
		break
	case "rejected", "cancelled":
		status = gateway.TransactionCancelled
		break
	case "refunded":
		status = gateway.TransactionRefunded
		break
	case "approved":
		status = gateway.TransactionPaid
		break
	default:
		return nil, nil
	}
	return &gateway.GetPaymentOutput{
		Id:             strconv.Itoa(paymentResponse.ID),
		QrCode:         paymentResponse.PointOfInteraction.TransactionData.QRCode,
		ExpiresAt:      paymentResponse.DateOfExpiration,
		PaymentGateway: mercadoPago,
		QrCodeLink:     paymentResponse.PointOfInteraction.TransactionData.TicketURL,
		Status:         status,
	}, nil
}

func (g *MercadoPagoPixGateway) RefundPix(ctx context.Context, input gateway.RefundPixInput) error {
	paymentId, err := strconv.Atoi(input.PaymentId)
	if err != nil {
		return err
	}
	cfg, err := g.getConfig(ctx, input.MerchantId)
	if err != nil {
		return err
	}
	client := refund.NewClient(cfg)
	_, err = client.Create(ctx, paymentId)
	return err
}
