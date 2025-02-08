package main

import (
	"bom-pedido-api/internal/application/factory"
	adminUseCase "bom-pedido-api/internal/application/usecase/admin"
	categoryUseCase "bom-pedido-api/internal/application/usecase/category"
	customerUseCase "bom-pedido-api/internal/application/usecase/customer"
	orderUseCase "bom-pedido-api/internal/application/usecase/order"
	productUseCase "bom-pedido-api/internal/application/usecase/product"
	shoppingCartUseCase "bom-pedido-api/internal/application/usecase/shopping_cart"
	"bom-pedido-api/internal/infra/config"
	"bom-pedido-api/internal/infra/handler"
	"bom-pedido-api/internal/infra/query"
	"bom-pedido-api/pkg/http"
	"bom-pedido-api/pkg/mongo"
)

func ConfigureRoutes(
	environment *config.Environment,
	applicationFactory *factory.ApplicationFactory,
	mongoDatabase *mongo.Database,
	server http.Server,
) {
	var (
		adminHandler = handler.NewAdminHandler(
			adminUseCase.NewSendAuthenticationMagicLink(environment.AdminMagicLinkBaseUrl, applicationFactory),
		)
		callbackHandler = handler.NewCallbackHandler(applicationFactory.EventHandler)
		categoryHandler = handler.NewCategoryHandler(
			categoryUseCase.NewCreateCategory(applicationFactory),
		)
		customerHandler = handler.NewCustomerHandler(
			customerUseCase.NewGetCustomer(applicationFactory),
			customerUseCase.NewGoogleAuthenticateCustomer(applicationFactory),
		)
		healthHandler = handler.NewHealthHandler(mongoDatabase)
		orderHandler  = handler.NewOrderHandler(
			orderUseCase.NewApproveOrder(applicationFactory),
			orderUseCase.NewCancelOrder(applicationFactory),
			orderUseCase.NewFinishOrder(applicationFactory),
			orderUseCase.NewCloneOrder(applicationFactory),
			orderUseCase.NewMarkOrderDelivering(applicationFactory),
			orderUseCase.NewMarkOrderInProgress(applicationFactory),
			orderUseCase.NewMarkOrderAwaitingWithdraw(applicationFactory),
			orderUseCase.NewMarkOrderAwaitingDelivery(applicationFactory),
			orderUseCase.NewRejectOrder(applicationFactory),
		)
		productHandler = handler.NewProductHandler(
			productUseCase.NewCreateProduct(applicationFactory),
			query.NewProductQuery(mongoDatabase),
		)
		shoppingCartHandler = handler.NewShoppingCartHandler(
			shoppingCartUseCase.NewAddItemToShoppingCart(applicationFactory),
			shoppingCartUseCase.NewCheckoutShoppingCart(applicationFactory),
			shoppingCartUseCase.NewGetShoppingCart(applicationFactory),
			shoppingCartUseCase.NewDeleteShoppingCart(applicationFactory),
			shoppingCartUseCase.NewDeleteShoppingCartItem(applicationFactory),
		)

		middlewares = http.NewMiddlewares(applicationFactory.Locker, applicationFactory.TokenManager)
	)

	server.AddMiddleware(middlewares.RecoverMiddleware)
	server.AddMiddleware(middlewares.RequestIDMiddleware)
	server.AddMiddleware(handler.ErrorHandlerMiddleware())
	server.AddMiddleware(middlewares.TelemetryMiddleware)
	server.AddMiddleware(middlewares.AuthenticateMiddleware())
	server.AddMiddleware(middlewares.MetricMiddleware)

	server.Get("/health", healthHandler.Health)

	// Admin
	server.Post("/v1/admin/auth", adminHandler.SendAuthenticationLink)

	// Category
	server.Post("/v1/categories", categoryHandler.CreateCategory, middlewares.OnlyAdminMiddleware)

	// Callback
	server.Post("/v1/payments/callback/WOOVI", callbackHandler.Woovi)
	server.Post("/v1/payments/callback/MERCADO_PAGO/:orderId", callbackHandler.MercadoPago)

	// Customer
	server.Get("/v1/customers/me", customerHandler.GetCustomer, middlewares.OnlyCustomerMiddleware)
	server.Post("/v1/auth/customer", customerHandler.AuthenticateCustomer)

	// Order
	ordersMiddlewares := []http.Middleware{
		middlewares.LockByRequestPath("id"),
		middlewares.OnlyAdminMiddleware,
	}
	server.Post("/v1/orders/{id}/approve", orderHandler.ApproveOrder, ordersMiddlewares...)
	server.Post("/v1/orders/{id}/cancel", orderHandler.CancelOrder, ordersMiddlewares...)
	server.Post("/v1/orders/{id}/finish", orderHandler.FinishOrder, ordersMiddlewares...)
	server.Post("/v1/orders/{id}/await-delivery", orderHandler.MarkOrderAwaitingDelivery, ordersMiddlewares...)
	server.Post("/v1/orders/{id}/await-withdraw", orderHandler.MarkOrderAwaitingWithdraw, ordersMiddlewares...)
	server.Post("/v1/orders/{id}/delivering", orderHandler.MarkOrderDelivering, ordersMiddlewares...)
	server.Post("/v1/orders/{id}/in-progress", orderHandler.MarkOrderInProgress, ordersMiddlewares...)
	server.Post("/v1/orders/{id}/reject", orderHandler.RejectOrder, ordersMiddlewares...)
	server.Post("/v1/orders/{id}/clone", orderHandler.CloneOrder, middlewares.OnlyCustomerMiddleware)

	// Product
	server.Get("/v1/products", productHandler.ListProducts)
	server.Post("/v1/products", productHandler.CreateProduct, middlewares.OnlyAdminMiddleware)

	// Shopping Cart
	server.Get("/v1/shopping-cart/me", shoppingCartHandler.GetShoppingCart, middlewares.OnlyCustomerMiddleware)
	server.Post("/v1/shopping-cart/checkout", shoppingCartHandler.Checkout, middlewares.OnlyCustomerMiddleware)
	server.Patch("/v1/shopping-cart/items", shoppingCartHandler.AddShoppingCartItem, middlewares.OnlyCustomerMiddleware)
	server.Delete("/v1/shopping-cart/items/:id", shoppingCartHandler.DeleteShoppingCartItem, middlewares.OnlyCustomerMiddleware)
	server.Delete("/v1/shopping-cart/me", shoppingCartHandler.DeleteShoppingCart, middlewares.OnlyCustomerMiddleware)
}
