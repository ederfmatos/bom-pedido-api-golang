package pix

import (
	"bom-pedido-api/application/gateway"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/infra/config"
	"bom-pedido-api/infra/http_client"
	"context"
	"fmt"
)

type MerchantConfigGatewayDecorator struct {
	gateways                               map[string]gateway.PixGateway
	merchantPaymentGatewayConfigRepository repository.MerchantPaymentGatewayConfigRepository
}

func NewMerchantConfigGatewayDecorator(
	environment config.PixPaymentGatewayEnv,
	merchantPaymentGatewayConfigRepository repository.MerchantPaymentGatewayConfigRepository,
) gateway.PixGateway {
	notificationUrl := environment.NotificationUrl
	expirationTimeInMinutes := environment.ExpirationTimeInMinutes
	wooviHttpClient := http_client.NewDefaultHttpClient(environment.WooviApiBaseUrl)
	return &MerchantConfigGatewayDecorator{
		merchantPaymentGatewayConfigRepository: merchantPaymentGatewayConfigRepository,
		gateways: map[string]gateway.PixGateway{
			mercadoPago: NewLogPixGatewayDecorator(NewMercadoPagoPixGateway(notificationUrl, expirationTimeInMinutes)),
			woovi:       NewLogPixGatewayDecorator(NewWooviPixGateway(wooviHttpClient, expirationTimeInMinutes)),
		},
	}
}

func (g *MerchantConfigGatewayDecorator) Name() string {
	return ""
}

func (g *MerchantConfigGatewayDecorator) CreateQrCodePix(ctx context.Context, input gateway.CreateQrCodePixInput) (*gateway.CreateQrCodePixOutput, error) {
	gatewayConfig, err := g.merchantPaymentGatewayConfigRepository.FindByMerchant(ctx, input.MerchantId)
	if err != nil {
		return nil, err
	}
	if gatewayConfig == nil {
		return nil, gateway.MerchantGatewayConfigNotFoundError
	}
	input.Credential = gatewayConfig.Credential
	paymentGateway, exists := g.gateways[gatewayConfig.PaymentGateway]
	if !exists {
		return nil, fmt.Errorf("payment gateway %s does not exists", gatewayConfig.PaymentGateway)
	}
	return paymentGateway.CreateQrCodePix(ctx, input)
}

func (g *MerchantConfigGatewayDecorator) GetPaymentById(ctx context.Context, input gateway.GetPaymentInput) (*gateway.GetPaymentOutput, error) {
	gatewayConfig, err := g.merchantPaymentGatewayConfigRepository.FindByMerchantAndGateway(ctx, input.MerchantId, input.PaymentGateway)
	if err != nil {
		return nil, err
	}
	if gatewayConfig == nil {
		return nil, gateway.MerchantGatewayConfigNotFoundError
	}
	input.Credential = gatewayConfig.Credential
	paymentGateway, exists := g.gateways[gatewayConfig.PaymentGateway]
	if !exists {
		return nil, fmt.Errorf("payment gateway %s does not exists", gatewayConfig.PaymentGateway)
	}
	return paymentGateway.GetPaymentById(ctx, input)
}

func (g *MerchantConfigGatewayDecorator) RefundPix(ctx context.Context, input gateway.RefundPixInput) error {
	gatewayConfig, err := g.merchantPaymentGatewayConfigRepository.FindByMerchantAndGateway(ctx, input.MerchantId, input.PaymentGateway)
	if err != nil {
		return err
	}
	if gatewayConfig == nil {
		return gateway.MerchantGatewayConfigNotFoundError
	}
	input.Credential = gatewayConfig.Credential
	paymentGateway, exists := g.gateways[gatewayConfig.PaymentGateway]
	if !exists {
		return fmt.Errorf("payment gateway %s does not exists", gatewayConfig.PaymentGateway)
	}
	return paymentGateway.RefundPix(ctx, input)
}
