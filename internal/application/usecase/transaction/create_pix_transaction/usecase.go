package create_pix_transaction

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/internal/application/lock"
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity"
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
	order, err := uc.orderRepository.FindById(ctx, input.OrderId)
	if err != nil || order == nil || !order.IsPixInApp() {
		return err
	}
	existsTransaction, err := uc.transactionRepository.ExistsByOrderId(ctx, order.Id)
	if err != nil || existsTransaction {
		return err
	}
	pixPayment, err := uc.pixGateway.GetPaymentById(ctx, gateway.GetPaymentInput{
		PaymentId:      input.PaymentId,
		MerchantId:     order.MerchantId,
		PaymentGateway: input.PaymentGateway,
	})
	if err != nil || pixPayment == nil {
		return err
	}
	pixTransaction := entity.NewPixTransaction(pixPayment.Id, order.Id, pixPayment.QrCode, pixPayment.PaymentGateway, pixPayment.QrCodeLink, order.Amount)
	err = uc.transactionRepository.CreatePixTransaction(ctx, pixTransaction)
	if err != nil {
		return err
	}
	return uc.eventEmitter.Emit(ctx, event.NewPixTransactionCreated(order.Id, pixTransaction.PaymentId))
}
