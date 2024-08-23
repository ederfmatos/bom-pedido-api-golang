package order

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/http/middlewares"
	"bom-pedido-api/infra/http/order/approve"
	"bom-pedido-api/infra/http/order/cancel"
	"bom-pedido-api/infra/http/order/clone"
	"bom-pedido-api/infra/http/order/finish"
	"bom-pedido-api/infra/http/order/mark_awaiting_delivery"
	"bom-pedido-api/infra/http/order/mark_awaiting_withdraw"
	"bom-pedido-api/infra/http/order/mark_delivering"
	"bom-pedido-api/infra/http/order/mark_in_progress"
	"bom-pedido-api/infra/http/order/reject"
	"github.com/labstack/echo/v4"
)

func ConfigureRoutes(server *echo.Group, applicationFactory *factory.ApplicationFactory) {
	orderRoutes := server.Group("/v1/orders/:id", middlewares.OnlyAdmin, middlewares.LockByParam("id", applicationFactory))
	orderRoutes.POST("/approve", approve.Handle(applicationFactory))
	orderRoutes.POST("/reject", reject.Handle(applicationFactory))
	orderRoutes.POST("/cancel", cancel.Handle(applicationFactory))
	orderRoutes.POST("/finish", finish.Handle(applicationFactory))
	orderRoutes.POST("/in-progress", mark_in_progress.Handle(applicationFactory))
	orderRoutes.POST("/delivering", mark_delivering.Handle(applicationFactory))
	orderRoutes.POST("/awaiting-withdraw", mark_awaiting_withdraw.Handle(applicationFactory))
	orderRoutes.POST("/awaiting-delivery", mark_awaiting_delivery.Handle(applicationFactory))
	orderRoutes.POST("/clone", clone.Handle(applicationFactory))
}
