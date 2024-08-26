package pix

import (
	"bom-pedido-api/application/gateway"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/errors"
	"bom-pedido-api/infra/telemetry"
	"context"
	"fmt"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
	"log/slog"
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

func (g *MercadoPagoPixGateway) CreateQrCodePix(ctx context.Context, input gateway.CreateQrCodePixInput) (*gateway.CreateQrCodePixOutput, error) {
	slog.Info("Iniciando criação de pagamento PIX no Mercado Pago")
	ctx, span := telemetry.StartSpan(ctx, "MercadoPagoPixGateway.CreateQrCodePix")
	defer span.End()
	gatewayConfig, err := g.merchantPaymentGatewayConfigRepository.FindByMerchantAndGateway(ctx, input.Merchant.Id, mercadoPago)
	defer func() {
		if err != nil {
			slog.Error("Ocorreu um erro na criação de pagamento Pix no Mercado Pago", "error", err)
		}
	}()
	if err != nil {
		return nil, err
	}
	if gatewayConfig == nil {
		return nil, errors.New("gateway config not found")
	}
	cfg, err := config.New(gatewayConfig.AccessToken, config.WithPlatformID(applicationName))
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
			FirstName: input.Merchant.Name,
			Email:     input.Merchant.Email,
		},
		Metadata: map[string]any{
			"orderId": input.InternalOrderId,
		},
	}
	_, mercadoPagoSpan := telemetry.StartSpan(ctx, "MercadoPagoPixGateway.CreatePayment")
	defer mercadoPagoSpan.End()
	resource, err := client.Create(ctx, request)
	if err != nil {
		mercadoPagoSpan.RecordError(err)
		return nil, err
	}
	slog.Info("Sucesso na criação de pagamento PIX no Mercado Pago")
	return &gateway.CreateQrCodePixOutput{
		Id:             strconv.Itoa(resource.ID),
		QrCode:         resource.PointOfInteraction.TransactionData.QRCode,
		ExpiresAt:      expiresAt,
		PaymentGateway: mercadoPago,
		QrCodeLink:     resource.PointOfInteraction.TransactionData.TicketURL,
	}, nil
}
