package cancel_pix_transaction

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/lock"
	"bom-pedido-api/internal/application/repository"
	"context"
	"time"
)

type (
	UseCase struct {
		orderRepository       repository.OrderRepository
		transactionRepository repository.TransactionRepository
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
		eventEmitter:          factory.EventEmitter,
		locker:                factory.Locker,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input Input) error {
	lockKey, err := uc.locker.Lock(ctx, time.Second*30, "CANCEL_PIX_TRANSACTION_", input.OrderId)
	if err != nil {
		return err
	}
	defer uc.locker.Release(ctx, lockKey)
	pixTransaction, err := uc.transactionRepository.FindByOrderId(ctx, input.OrderId)
	if err != nil || pixTransaction == nil || !pixTransaction.IsCreated() {
		return err
	}
	pixTransaction.Cancel()
	if err = uc.transactionRepository.UpdatePixTransaction(ctx, pixTransaction); err != nil {
		return err
	}
	return uc.eventEmitter.Emit(ctx, event.NewPixTransactionCancelled(pixTransaction.OrderId, pixTransaction.Id))
}
