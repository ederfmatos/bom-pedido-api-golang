package pay_pix_transaction

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/gateway"
	"bom-pedido-api/application/repository"
	"context"
)

type (
	UseCase struct {
		orderRepository       repository.OrderRepository
		transactionRepository repository.TransactionRepository
		merchantRepository    repository.MerchantRepository
		pixGateway            gateway.PixGateway
		eventEmitter          event.Emitter
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
	}
}

func (uc *UseCase) Execute(ctx context.Context, input Input) error {
	anOrder, err := uc.orderRepository.FindById(ctx, input.OrderId)
	if err != nil || anOrder == nil || !anOrder.IsPixInApp() || !anOrder.IsAwaitingPayment() {
		return err
	}
	pendingTransaction, err := uc.transactionRepository.FindPendingByOrderId(ctx, anOrder.Id)
	if err != nil || pendingTransaction == nil {
		return err
	}
	paymentStatus, err := uc.pixGateway.GetPaymentStatus(ctx, anOrder.MerchantId, pendingTransaction.Id)
	if err != nil || paymentStatus != nil {
		return err
	}
	if *paymentStatus != gateway.TransactionPaid {
		return err
	}
	pendingTransaction.Pay()
	err = uc.transactionRepository.UpdatePixTransaction(ctx, pendingTransaction)
	if err != nil {
		return err
	}
	return uc.eventEmitter.Emit(ctx, event.NewTransactionPaid(anOrder.Id, pendingTransaction.Id))
}
