package messaging

import (
	"bom-pedido-api/application/factory"
)

func HandleEvents(factory *factory.ApplicationFactory) {
	HandleShoppingCart(factory)
	HandleOrderEvents(factory)
	HandleTransactionEvents(factory)
	HandleEmailEvents(factory)
	HandlePaymentEvents(factory)
	HandleNotificationEvents(factory)
	HandleTransactionCallback(factory)
}
