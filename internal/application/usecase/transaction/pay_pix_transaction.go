package transaction

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/internal/application/lock"
	"bom-pedido-api/internal/application/repository"
	"context"
)

type (
	PayPixTransactionUseCase struct {
		orderRepository       repository.OrderRepository
		transactionRepository repository.TransactionRepository
		pixGateway            gateway.PixGateway
		eventEmitter          event.Emitter
		locker                lock.Locker
	}

	PayPixTransactionInput struct {
		OrderId string
	}
)

func NewPayPixTransaction(factory *factory.ApplicationFactory) *PayPixTransactionUseCase {
	return &PayPixTransactionUseCase{
		orderRepository:       factory.OrderRepository,
		transactionRepository: factory.TransactionRepository,
		pixGateway:            factory.PixGateway,
		eventEmitter:          factory.EventEmitter,
		locker:                factory.Locker,
	}
}

func (uc *PayPixTransactionUseCase) Execute(ctx context.Context, input PayPixTransactionInput) error {
	lockKey, err := uc.locker.Lock(ctx, "PAY_PIX_TRANSACTION_", input.OrderId)
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
	payment, err := uc.pixGateway.GetPaymentById(ctx, gateway.GetPaymentInput{
		PaymentId:      pixTransaction.PaymentId,
		MerchantId:     order.MerchantId,
		PaymentGateway: pixTransaction.PaymentGateway,
	})
	if err != nil || payment == nil || payment.Status != gateway.TransactionPaid {
		return err
	}
	pixTransaction.Pay()
	if err = uc.transactionRepository.UpdatePixTransaction(ctx, pixTransaction); err != nil {
		return err
	}
	return uc.eventEmitter.Emit(ctx, event.NewPixTransactionPaid(order.Id, pixTransaction.Id))
}
