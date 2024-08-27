package pay_pix_transaction

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/gateway"
	"bom-pedido-api/application/lock"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/transaction"
	"context"
	"errors"
	"time"
)

type (
	UseCase struct {
		orderRepository       repository.OrderRepository
		transactionRepository repository.TransactionRepository
		merchantRepository    repository.MerchantRepository
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
		merchantRepository:    factory.MerchantRepository,
		pixGateway:            factory.PixGateway,
		eventEmitter:          factory.EventEmitter,
		locker:                factory.Locker,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input Input) error {
	lockKey, err := uc.locker.Lock(ctx, time.Second*30, "PAY_PIX_TRANSACTION_", input.OrderId)
	if err != nil {
		return err
	}
	defer uc.locker.Release(ctx, lockKey)
	anOrder, err := uc.orderRepository.FindById(ctx, input.OrderId)
	if err != nil || anOrder == nil || !anOrder.IsPixInApp() || !anOrder.IsAwaitingPayment() {
		return err
	}
	pixTransaction, err := uc.transactionRepository.FindByOrderId(ctx, anOrder.Id)
	if err != nil || pixTransaction == nil {
		return err
	}
	paymentStatus, err := uc.pixGateway.GetPaymentStatus(ctx, anOrder.MerchantId, pixTransaction.Id)
	if err != nil || paymentStatus != nil {
		return err
	}
	if *paymentStatus != gateway.TransactionPaid {
		return err
	}
	err = pixTransaction.Pay()
	if errors.Is(err, transaction.AlreadyPaidError) {
		return nil
	}
	err = uc.transactionRepository.UpdatePixTransaction(ctx, pixTransaction)
	if err != nil {
		return err
	}
	return uc.eventEmitter.Emit(ctx, event.NewTransactionPaid(anOrder.Id, pixTransaction.Id))
}
