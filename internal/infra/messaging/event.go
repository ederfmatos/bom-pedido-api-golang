package messaging

import (
	"bom-pedido-api/internal/application/factory"
)

func HandleEvents(factory *factory.ApplicationFactory) {
	HandleShoppingCart(factory)
	HandleOrderEvents(factory)
	HandleTransactionEvents(factory)
	HandleEmailEvents(factory)
	HandlePaymentEvents(factory)
	HandleNotificationEvents(factory)
	HandleTransactiOnback(factory)
}
