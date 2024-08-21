package shopping_cart

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/http/middlewares"
	"bom-pedido-api/infra/http/shopping_cart/add_item_to_shopping_cart"
	"bom-pedido-api/infra/http/shopping_cart/checkout_shopping_cart"
	"github.com/labstack/echo/v4"
)

func ConfigureRoutes(server *echo.Group, applicationFactory *factory.ApplicationFactory) {
	shoppingCartRoutes := server.Group("/v1/shopping-cart", middlewares.LockByCustomerId(applicationFactory))
	shoppingCartRoutes.POST("/checkout", checkout_shopping_cart.Handle(applicationFactory))
	shoppingCartRoutes.PATCH("/items", add_item_to_shopping_cart.Handle(applicationFactory))
}
