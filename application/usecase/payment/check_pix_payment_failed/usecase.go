package check_pix_payment_failed

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/gateway"
	"bom-pedido-api/application/lock"
	"bom-pedido-api/application/repository"
	"context"
	"time"
)

type (
	UseCase struct {
		orderRepository       repository.OrderRepository
		transactionRepository repository.TransactionRepository
		pixGateway            gateway.PixGateway
		eventEmitter          event.Emitter
		locker                lock.Locker
	}

	Input struct {
		OrderId string
	}
)

func New(factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{
		orderRepository:       factory.OrderRepository,
		transactionRepository: factory.TransactionRepository,
		pixGateway:            factory.PixGateway,
		eventEmitter:          factory.EventEmitter,
		locker:                factory.Locker,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input Input) error {
	lockKey, err := uc.locker.Lock(ctx, time.Second*30, "CHECK_PIX_PAYMENT_FAILED_", input.OrderId)
	if err != nil {
		return err
	}
	defer uc.locker.Release(ctx, lockKey)
	anOrder, err := uc.orderRepository.FindById(ctx, input.OrderId)
	if err != nil || anOrder == nil || !anOrder.IsPixInApp() || !anOrder.IsAwaitingPayment() {
		return err
	}
	pixTransaction, err := uc.transactionRepository.FindByOrderId(ctx, anOrder.Id)
	if err != nil || pixTransaction == nil || pixTransaction.IsPaid() {
		return err
	}
	getPaymentInput := gateway.GetPaymentInput{
		MerchantId:     anOrder.MerchantId,
		PaymentId:      pixTransaction.PaymentId,
		PaymentGateway: pixTransaction.PaymentGateway,
	}
	pixPayment, err := uc.pixGateway.GetPaymentById(ctx, getPaymentInput)
	if err != nil || pixPayment == nil || pixPayment.Status != gateway.TransactionCancelled {
		return err
	}
	return uc.eventEmitter.Emit(ctx, event.NewPixPaymentCancelled(anOrder.Id, pixPayment.Id, pixPayment.PaymentGateway))
}
