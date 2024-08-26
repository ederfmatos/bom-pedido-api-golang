package messaging

import (
	"bom-pedido-api/application/factory"
)

func HandleEvents(factory *factory.ApplicationFactory) {
	HandleShoppingCart(factory)
	HandleOrderEvents(factory)
	HandleEmailEvents(factory)
	HandleTransactionCallback(factory)
}
