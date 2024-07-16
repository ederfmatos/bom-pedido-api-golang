package http

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/presentation/http/customer/get_customer"
	"bom-pedido-api/presentation/http/customer/google_auth_customer"
	"bom-pedido-api/presentation/http/health"
	"bom-pedido-api/presentation/http/middlewares"
	"bom-pedido-api/presentation/http/order/approve"
	"bom-pedido-api/presentation/http/order/cancel"
	"bom-pedido-api/presentation/http/order/finish"
	"bom-pedido-api/presentation/http/order/mark_awaiting_delivery"
	"bom-pedido-api/presentation/http/order/mark_awaiting_withdraw"
	"bom-pedido-api/presentation/http/order/mark_delivering"
	"bom-pedido-api/presentation/http/order/mark_in_progress"
	"bom-pedido-api/presentation/http/order/reject"
	"bom-pedido-api/presentation/http/products/create_product"
	"bom-pedido-api/presentation/http/products/list_products"
	"bom-pedido-api/presentation/http/shopping_cart/add_item_to_shopping_cart"
	"bom-pedido-api/presentation/http/shopping_cart/checkout_shopping_cart"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Server(applicationFactory *factory.ApplicationFactory) *echo.Echo {
	server := echo.New()
	server.Use(middleware.Recover())
	server.Use(middleware.RequestID())
	server.Use(middlewares.AuthenticateMiddleware(applicationFactory))
	server.HTTPErrorHandler = middlewares.HandleError
	SetRoutes(server, applicationFactory)
	return server
}

func SetRoutes(server *echo.Echo, applicationFactory *factory.ApplicationFactory) {
	api := server.Group("/api")
	shoppingCartRoutes := api.Group("/v1/shopping-cart", middlewares.LockByCustomerId(applicationFactory))
	shoppingCartRoutes.POST("/checkout", checkout_shopping_cart.Handle(applicationFactory))
	shoppingCartRoutes.PATCH("/items", add_item_to_shopping_cart.Handle(applicationFactory))

	api.POST("/v1/products", create_product.Handle(applicationFactory))
	api.GET("/v1/products", list_products.Handle(applicationFactory))
	api.POST("/v1/auth/google/customer", google_auth_customer.Handle(applicationFactory))
	api.GET("/v1/customers/me", get_customer.Handle(applicationFactory))

	orderRoutes := api.Group("/v1/orders/:id", middlewares.LockByParam("id", applicationFactory))
	orderRoutes.POST("/approve_order", approve.Handle(applicationFactory))
	orderRoutes.POST("/reject_order", reject.Handle(applicationFactory))
	orderRoutes.POST("/cancel_order", cancel.Handle(applicationFactory))
	orderRoutes.POST("/finish_order", finish.Handle(applicationFactory))
	orderRoutes.POST("/in-progress", mark_in_progress.Handle(applicationFactory))
	orderRoutes.POST("/delivering", mark_delivering.Handle(applicationFactory))
	orderRoutes.POST("/awaiting-withdraw", mark_awaiting_withdraw.Handle(applicationFactory))
	orderRoutes.POST("/awaiting-delivery", mark_awaiting_delivery.Handle(applicationFactory))

	api.GET("/health", health.Handle)
}
