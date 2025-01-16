package payment

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/internal/application/lock"
	"bom-pedido-api/internal/application/repository"
	"context"
	"time"
)

type (
	RefundPixPaymentUseCase struct {
		orderRepository       repository.OrderRepository
		transactionRepository repository.TransactionRepository
		pixGateway            gateway.PixGateway
		eventEmitter          event.Emitter
		locker                lock.Locker
	}

	RefundPixPaymentInput struct {
		OrderId string
	}
)

func NewRefundPixPayment(factory *factory.ApplicationFactory) *RefundPixPaymentUseCase {
	return &RefundPixPaymentUseCase{
		orderRepository:       factory.OrderRepository,
		transactionRepository: factory.TransactionRepository,
		pixGateway:            factory.PixGateway,
		eventEmitter:          factory.EventEmitter,
		locker:                factory.Locker,
	}
}

func (uc *RefundPixPaymentUseCase) Execute(ctx context.Context, input RefundPixPaymentInput) error {
	lockKey, err := uc.locker.Lock(ctx, time.Second*30, "REFUND_PIX_PAYMENT_", input.OrderId)
	if err != nil {
		return err
	}
	defer uc.locker.Release(ctx, lockKey)
	order, err := uc.orderRepository.FindById(ctx, input.OrderId)
	if err != nil || order == nil || !order.IsPixInApp() || order.IsAwaitingPayment() {
		return err
	}
	pixTransaction, err := uc.transactionRepository.FindByOrderId(ctx, order.Id)
	if err != nil || pixTransaction == nil {
		return err
	}
	payment, err := uc.pixGateway.GetPaymentById(ctx, gateway.GetPaymentInput{
		PaymentId:      pixTransaction.PaymentId,
		MerchantId:     order.MerchantId,
		PaymentGateway: pixTransaction.PaymentGateway,
	})
	if err != nil || payment == nil || payment.Status != gateway.TransactionPaid {
		return err
	}
	refundInput := gateway.RefundPixInput{
		PaymentId:      pixTransaction.PaymentId,
		MerchantId:     order.MerchantId,
		Amount:         order.Amount,
		PaymentGateway: pixTransaction.PaymentGateway,
	}
	if err = uc.pixGateway.RefundPix(ctx, refundInput); err != nil {
		return err
	}
	return uc.eventEmitter.Emit(ctx, event.NewPixPaymentRefunded(order.Id, pixTransaction.PaymentId))
}
