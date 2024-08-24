package messaging

import (
	"bom-pedido-api/application/factory"
)

func HandleEvents(factory *factory.ApplicationFactory) {
	HandleShoppingCart(factory)
	HandleProductEvents(factory)
	HandleOrderEvents(factory)
	HandleEmailEvents(factory)
}
