package pix

import (
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/internal/infra/telemetry"
	"context"
	"log/slog"
)

type LogPixGatewayDecorator struct {
	delegate gateway.PixGateway
	logger   *slog.Logger
}

func NewLogPixGatewayDecorator(delegate gateway.PixGateway) gateway.PixGateway {
	return &LogPixGatewayDecorator{
		delegate: delegate,
		logger:   slog.With("paymentGateway", delegate.Name()),
	}
}

func (g *LogPixGatewayDecorator) Name() string {
	return g.delegate.Name()
}

func (g *LogPixGatewayDecorator) CreateQrCodePix(ctx context.Context, input gateway.CreateQrCodePixInput) (*gateway.CreateQrCodePixOutput, error) {
	g.logger.Info("Iniciando criação de pagamento PIX")
	ctx, span := telemetry.StartSpan(ctx, "PixGateway.CreateQrCodePix")
	defer span.End()
	output, err := g.delegate.CreateQrCodePix(ctx, input)
	if err != nil {
		g.logger.Error("Ocorreu um erro na criação de pagamento Pix", "error", err)
		span.RecordError(err)
		return nil, err
	}
	g.logger.Info("Sucesso na criação de pagamento PIX")
	return output, nil
}

func (g *LogPixGatewayDecorator) GetPaymentById(ctx context.Context, input gateway.GetPaymentInput) (*gateway.GetPaymentOutput, error) {
	g.logger.Info("Iniciando busca de pagamento PIX")
	ctx, span := telemetry.StartSpan(ctx, "PixGateway.GetPaymentById")
	defer span.End()
	output, err := g.delegate.GetPaymentById(ctx, input)
	if err != nil {
		g.logger.Error("Ocorreu um erro na busca de pagamento Pix", "error", err)
		return nil, err
	}
	g.logger.Info("Sucesso na busca de pagamento PIX")
	return output, nil
}

func (g *LogPixGatewayDecorator) RefundPix(ctx context.Context, input gateway.RefundPixInput) error {
	g.logger.Info("Iniciando estorno de pagamento PIX")
	ctx, span := telemetry.StartSpan(ctx, "PixGateway.RefundPix")
	defer span.End()
	err := g.delegate.RefundPix(ctx, input)
	if err != nil {
		g.logger.Error("Ocorreu um erro na estorno de pagamento Pix", "error", err)
		return err
	}
	g.logger.Info("Sucesso na estorno de pagamento PIX")
	return nil
}
