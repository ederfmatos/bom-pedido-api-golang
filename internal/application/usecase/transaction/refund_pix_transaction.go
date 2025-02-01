package transaction

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/internal/application/lock"
	"bom-pedido-api/internal/application/repository"
	"context"
)

type (
	RefundPixTransactionUseCase struct {
		orderRepository       repository.OrderRepository
		transactionRepository repository.TransactionRepository
		pixGateway            gateway.PixGateway
		eventEmitter          event.Emitter
		locker                lock.Locker
	}

	RefundPixTransactionInput struct {
		OrderId string
	}
)

func NewRefundPixTransaction(
	orderRepository repository.OrderRepository,
	transactionRepository repository.TransactionRepository,
	pixGateway gateway.PixGateway,
	eventEmitter event.Emitter,
	locker lock.Locker,
) *RefundPixTransactionUseCase {
	return &RefundPixTransactionUseCase{
		orderRepository:       orderRepository,
		transactionRepository: transactionRepository,
		pixGateway:            pixGateway,
		eventEmitter:          eventEmitter,
		locker:                locker,
	}
}

func (uc *RefundPixTransactionUseCase) Execute(ctx context.Context, input RefundPixTransactionInput) error {
	lockKey, err := uc.locker.Lock(ctx, "REFUND_PIX_TRANSACTION_", input.OrderId)
	if err != nil {
		return err
	}
	defer uc.locker.Release(ctx, lockKey)
	order, err := uc.orderRepository.FindById(ctx, input.OrderId)
	if err != nil || order == nil || !order.IsPixInApp() || order.IsAwaitingPayment() {
		return err
	}
	pixTransaction, err := uc.transactionRepository.FindByOrderId(ctx, order.Id)
	if err != nil || pixTransaction == nil || !pixTransaction.IsPaid() {
		return err
	}
	payment, err := uc.pixGateway.GetPaymentById(ctx, gateway.GetPaymentInput{
		PaymentId:      pixTransaction.PaymentId,
		MerchantId:     order.MerchantId,
		PaymentGateway: pixTransaction.PaymentGateway,
	})
	if err != nil || payment == nil || payment.Status != gateway.TransactionRefunded {
		return err
	}
	pixTransaction.Refund()
	if err = uc.transactionRepository.UpdatePixTransaction(ctx, pixTransaction); err != nil {
		return err
	}
	return uc.eventEmitter.Emit(ctx, event.NewPixTransactionRefunded(order.Id, pixTransaction.Id))
}
