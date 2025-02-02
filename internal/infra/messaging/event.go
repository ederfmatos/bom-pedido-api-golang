package messaging

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/usecase/notification"
	"bom-pedido-api/internal/application/usecase/order"
	"bom-pedido-api/internal/application/usecase/payment"
	"bom-pedido-api/internal/application/usecase/shopping_cart"
	"bom-pedido-api/internal/application/usecase/transaction"
)

func HandleEvents(factory *factory.ApplicationFactory) {
	var (
		eventHandler = factory.EventHandler

		emailHandler = NewEmailHandler(factory.EmailGateway)

		notificationHandler = NewNotificationHandler(
			notification.NewNotifyCustomerOrderStatusChanged(factory),
			notification.NewSendNotification(factory),
		)

		orderHandler = NewOrderHandler(
			order.NewSaveOrderHistory(factory),
			order.NewAwaitApprovalOrderUseCase(factory),
			order.NewFailOrderPayment(factory),
		)

		paymentHandler = NewPaymentHandler(
			payment.NewCheckPixPaymentFailed(factory),
			payment.NewCreatePixPayment(factory),
			payment.NewRefundPixPayment(factory),
		)

		shoppingCartHandler = NewShoppingCartHandlerHandler(
			shopping_cart.NewDeleteShoppingCart(factory),
		)

		transactionHandler = NewTransactionHandler(
			transaction.NewPayPixTransaction(factory),
			transaction.NewCreatePixTransaction(factory),
			transaction.NewRefundPixTransaction(
				factory.OrderRepository,
				factory.TransactionRepository,
				factory.PixGateway,
				factory.EventEmitter,
				factory.Locker,
			),
			transaction.NewCancelPixTransaction(
				factory.OrderRepository,
				factory.TransactionRepository,
				factory.EventEmitter,
				factory.Locker,
			),
		)
	)

	// Email
	eventHandler.OnEvent("SEND_EMAIL", emailHandler.HandleSendEmail)

	// Notification
	eventHandler.OnEvent("NOTIFY_CUSTOMER_ORDER_STATUS_CHANGED", notificationHandler.NotifyCustomerOrderStatusChanged)

	// Order
	eventHandler.OnEvent("AWAIT_APPROVAL_ORDER", orderHandler.HandleAwaitApprovalOrder)
	eventHandler.OnEvent("ORDER_PAYMENT_FAILED", orderHandler.HandleOrderPaymentFailed)
	eventHandler.OnEvent("SAVE_ORDER_STATUS_HISTORY", orderHandler.HandleOrderStatusChanged)

	// Payment
	eventHandler.OnEvent("CREATE_PIX_PAYMENT", paymentHandler.CreatePixPayment)
	eventHandler.OnEvent("REFUND_PIX_PAYMENT", paymentHandler.RefundPixPayment)
	eventHandler.OnEvent("CHECK_PIX_PAYMENT_FAILED", paymentHandler.CheckPixPaymentFailed)

	// Shopping Cart
	eventHandler.OnEvent("ORDER_CREATED", shoppingCartHandler.DeleteShoppingCart)
	eventHandler.OnEvent("DELETE_SHOPPING_CART", shoppingCartHandler.DeleteShoppingCart)

	// Transaction
	eventHandler.OnEvent("PAY_PIX_TRANSACTION", transactionHandler.HandlePayPixTransaction)
	eventHandler.OnEvent("CREATE_PIX_TRANSACTION", transactionHandler.HandleCreatePixTransaction)
	eventHandler.OnEvent("REFUND_PIX_TRANSACTION", transactionHandler.HandleRefundPixTransaction)
	eventHandler.OnEvent("CANCEL_PIX_TRANSACTION", transactionHandler.HandleCancelPixTransaction)
}
