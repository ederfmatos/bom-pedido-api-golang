package pix

import (
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/internal/application/repository"
	"context"
	"fmt"
)

type MerchantPixGatewayMediator struct {
	gateways                               map[string]gateway.PixGateway
	merchantPaymentGatewayConfigRepository repository.MerchantPaymentGatewayConfigRepository
}

func NewMerchantPixGatewayMediator(
	merchantPaymentGatewayConfigRepository repository.MerchantPaymentGatewayConfigRepository,
	gateways map[string]gateway.PixGateway,
) gateway.PixGateway {
	return &MerchantPixGatewayMediator{
		merchantPaymentGatewayConfigRepository: merchantPaymentGatewayConfigRepository,
		gateways:                               gateways,
	}
}

func (g *MerchantPixGatewayMediator) Name() string {
	return "MEDIATOR"
}

func (g *MerchantPixGatewayMediator) CreateQrCodePix(ctx context.Context, input gateway.CreateQrCodePixInput) (*gateway.CreateQrCodePixOutput, error) {
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

func (g *MerchantPixGatewayMediator) GetPaymentById(ctx context.Context, input gateway.GetPaymentInput) (*gateway.GetPaymentOutput, error) {
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

func (g *MerchantPixGatewayMediator) RefundPix(ctx context.Context, input gateway.RefundPixInput) error {
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
