package pix

import (
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/internal/infra/telemetry"
	"bom-pedido-api/pkg/log"
	"context"
)

type LogPixGatewayDecorator struct {
	delegate gateway.PixGateway
}

func NewLogPixGatewayDecorator(delegate gateway.PixGateway) gateway.PixGateway {
	return &LogPixGatewayDecorator{delegate: delegate}
}

func (g *LogPixGatewayDecorator) Name() string {
	return g.delegate.Name()
}

func (g *LogPixGatewayDecorator) CreateQrCodePix(ctx context.Context, input gateway.CreateQrCodePixInput) (*gateway.CreateQrCodePixOutput, error) {
	log.Info("Iniciando criação de pagamento PIX", "paymentGateway", g.delegate.Name())
	ctx, span := telemetry.StartSpan(ctx, "PixGateway.CreateQrCodePix")
	defer span.End()
	output, err := g.delegate.CreateQrCodePix(ctx, input)
	if err != nil {
		log.Error("Ocorreu um erro na criação de pagamento Pix", err, "paymentGateway", g.delegate.Name())
		span.RecordError(err)
		return nil, err
	}
	log.Info("Sucesso na criação de pagamento PIX", "paymentGateway", g.delegate.Name())
	return output, nil
}

func (g *LogPixGatewayDecorator) GetPaymentById(ctx context.Context, input gateway.GetPaymentInput) (*gateway.GetPaymentOutput, error) {
	log.Info("Iniciando busca de pagamento PIX", "paymentGateway", g.delegate.Name())
	ctx, span := telemetry.StartSpan(ctx, "PixGateway.GetPaymentById")
	defer span.End()
	output, err := g.delegate.GetPaymentById(ctx, input)
	if err != nil {
		log.Error("Ocorreu um erro na busca de pagamento Pix", err, "paymentGateway", g.delegate.Name())
		return nil, err
	}
	log.Info("Sucesso na busca de pagamento PIX", "paymentGateway", g.delegate.Name())
	return output, nil
}

func (g *LogPixGatewayDecorator) RefundPix(ctx context.Context, input gateway.RefundPixInput) error {
	log.Info("Iniciando estorno de pagamento PIX", "paymentGateway", g.delegate.Name())
	ctx, span := telemetry.StartSpan(ctx, "PixGateway.RefundPix")
	defer span.End()
	err := g.delegate.RefundPix(ctx, input)
	if err != nil {
		log.Error("Ocorreu um erro na estorno de pagamento Pix", err, "paymentGateway", g.delegate.Name())
		return err
	}
	log.Info("Sucesso na estorno de pagamento PIX", "paymentGateway", g.delegate.Name())
	return nil
}
