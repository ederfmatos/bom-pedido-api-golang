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
	CheckPixPaymentFailedUseCase struct {
		orderRepository       repository.OrderRepository
		transactionRepository repository.TransactionRepository
		pixGateway            gateway.PixGateway
		eventEmitter          event.Emitter
		locker                lock.Locker
	}

	CheckPixPaymentFailedInput struct {
		OrderId string
	}
)

func NewCheckPixPaymentFailed(factory *factory.ApplicationFactory) *CheckPixPaymentFailedUseCase {
	return &CheckPixPaymentFailedUseCase{
		orderRepository:       factory.OrderRepository,
		transactionRepository: factory.TransactionRepository,
		pixGateway:            factory.PixGateway,
		eventEmitter:          factory.EventEmitter,
		locker:                factory.Locker,
	}
}

func (uc *CheckPixPaymentFailedUseCase) Execute(ctx context.Context, input CheckPixPaymentFailedInput) error {
	lockKey, err := uc.locker.Lock(ctx, time.Second*30, "CHECK_PIX_PAYMENT_FAILED_", input.OrderId)
	if err != nil {
		return err
	}
	defer uc.locker.Release(ctx, lockKey)
	order, err := uc.orderRepository.FindById(ctx, input.OrderId)
	if err != nil || order == nil || !order.IsPixInApp() || !order.IsAwaitingPayment() {
		return err
	}
	pixTransaction, err := uc.transactionRepository.FindByOrderId(ctx, order.Id)
	if err != nil || pixTransaction == nil || pixTransaction.IsPaid() {
		return err
	}
	getPaymentInput := gateway.GetPaymentInput{
		MerchantId:     order.MerchantId,
		PaymentId:      pixTransaction.PaymentId,
		PaymentGateway: pixTransaction.PaymentGateway,
	}
	pixPayment, err := uc.pixGateway.GetPaymentById(ctx, getPaymentInput)
	if err != nil || pixPayment == nil || pixPayment.Status != gateway.TransactionCancelled {
		return err
	}
	return uc.eventEmitter.Emit(ctx, event.NewPixPaymentCancelled(order.Id, pixPayment.Id, pixPayment.PaymentGateway))
}
