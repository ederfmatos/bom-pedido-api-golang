package refund_pix_payment

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
	lockKey, err := uc.locker.Lock(ctx, time.Second*30, "REFUND_PIX_PAYMENT_", input.OrderId)
	if err != nil {
		return err
	}
	defer uc.locker.Release(ctx, lockKey)
	anOrder, err := uc.orderRepository.FindById(ctx, input.OrderId)
	if err != nil || anOrder == nil || !anOrder.IsPixInApp() || anOrder.IsAwaitingPayment() {
		return err
	}
	pixTransaction, err := uc.transactionRepository.FindByOrderId(ctx, anOrder.Id)
	if err != nil || pixTransaction == nil {
		return err
	}
	payment, err := uc.pixGateway.GetPaymentById(ctx, gateway.GetPaymentInput{
		PaymentId:      pixTransaction.PaymentId,
		MerchantId:     anOrder.MerchantId,
		PaymentGateway: pixTransaction.PaymentGateway,
	})
	if err != nil || payment == nil || payment.Status != gateway.TransactionPaid {
		return err
	}
	refundInput := gateway.RefundPixInput{
		PaymentId:      pixTransaction.PaymentId,
		MerchantId:     anOrder.MerchantId,
		Amount:         anOrder.Amount,
		PaymentGateway: pixTransaction.PaymentGateway,
	}
	if err = uc.pixGateway.RefundPix(ctx, refundInput); err != nil {
		return err
	}
	return uc.eventEmitter.Emit(ctx, event.NewPixPaymentRefunded(anOrder.Id, pixTransaction.PaymentId))
}
