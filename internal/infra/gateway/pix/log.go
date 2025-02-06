package pix

import (
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/pkg/telemetry"
	"context"
)

type TelemetryPixGatewayDecorator struct {
	delegate gateway.PixGateway
}

func NewTelemetryPixGatewayDecorator(delegate gateway.PixGateway) gateway.PixGateway {
	return &TelemetryPixGatewayDecorator{delegate: delegate}
}

func (g *TelemetryPixGatewayDecorator) Name() string {
	return g.delegate.Name()
}

func (g *TelemetryPixGatewayDecorator) CreateQrCodePix(ctx context.Context, input gateway.CreateQrCodePixInput) (*gateway.CreateQrCodePixOutput, error) {
	return telemetry.StartSpan[*gateway.CreateQrCodePixOutput](ctx, "PixGateway.CreateQrCodePix", func(ctx context.Context) (*gateway.CreateQrCodePixOutput, error) {
		return g.delegate.CreateQrCodePix(ctx, input)
	})
}

func (g *TelemetryPixGatewayDecorator) GetPaymentById(ctx context.Context, input gateway.GetPaymentInput) (*gateway.GetPaymentOutput, error) {
	return telemetry.StartSpan[*gateway.GetPaymentOutput](ctx, "PixGateway.GetPaymentById", func(ctx context.Context) (*gateway.GetPaymentOutput, error) {
		return g.delegate.GetPaymentById(ctx, input)
	})
}

func (g *TelemetryPixGatewayDecorator) RefundPix(ctx context.Context, input gateway.RefundPixInput) error {
	return telemetry.StartSpanReturningError(ctx, "PixGateway.RefundPix", func(ctx context.Context) error {
		return g.delegate.RefundPix(ctx, input)
	})
}
