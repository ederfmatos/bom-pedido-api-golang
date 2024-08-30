package create_pix_transaction

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/gateway"
	"bom-pedido-api/application/lock"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/transaction"
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
		OrderId        string
		PaymentId      string
		PaymentGateway string
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
	lockKey, err := uc.locker.Lock(ctx, time.Second*30, "CREATE_PIX_TRANSACTION_", input.OrderId)
	if err != nil {
		return err
	}
	defer uc.locker.Release(ctx, lockKey)
	anOrder, err := uc.orderRepository.FindById(ctx, input.OrderId)
	if err != nil || anOrder == nil || !anOrder.IsPixInApp() {
		return err
	}
	existsTransaction, err := uc.transactionRepository.ExistsByOrderId(ctx, anOrder.Id)
	if err != nil || existsTransaction {
		return err
	}
	pixPayment, err := uc.pixGateway.GetPaymentById(ctx, gateway.GetPaymentInput{
		PaymentId:      input.PaymentId,
		MerchantId:     anOrder.MerchantId,
		PaymentGateway: input.PaymentGateway,
	})
	if err != nil || pixPayment == nil {
		return err
	}
	pixTransaction := transaction.NewPixTransaction(pixPayment.Id, anOrder.Id, pixPayment.QrCode, pixPayment.PaymentGateway, pixPayment.QrCodeLink, anOrder.Amount)
	err = uc.transactionRepository.CreatePixTransaction(ctx, pixTransaction)
	if err != nil {
		return err
	}
	return uc.eventEmitter.Emit(ctx, event.NewPixTransactionCreated(anOrder.Id, pixTransaction.PaymentId))
}
