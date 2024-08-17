package http

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/infra/config"
	"bom-pedido-api/infra/http/customer/get_customer"
	"bom-pedido-api/infra/http/customer/google_auth_customer"
	"bom-pedido-api/infra/http/health"
	"bom-pedido-api/infra/http/middlewares"
	"bom-pedido-api/infra/http/order/approve"
	"bom-pedido-api/infra/http/order/cancel"
	"bom-pedido-api/infra/http/order/finish"
	"bom-pedido-api/infra/http/order/mark_awaiting_delivery"
	"bom-pedido-api/infra/http/order/mark_awaiting_withdraw"
	"bom-pedido-api/infra/http/order/mark_delivering"
	"bom-pedido-api/infra/http/order/mark_in_progress"
	"bom-pedido-api/infra/http/order/reject"
	"bom-pedido-api/infra/http/products/create_product"
	"bom-pedido-api/infra/http/products/list_products"
	"bom-pedido-api/infra/http/shopping_cart/add_item_to_shopping_cart"
	"bom-pedido-api/infra/http/shopping_cart/checkout_shopping_cart"
	"context"
	"database/sql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	server         *echo.Echo
	mongoClient    *mongo.Client
	database       *sql.DB
	redisClient    *redis.Client
	tracerProvider *trace.TracerProvider
	environment    *config.Environment
}

func NewServer(database *sql.DB, redisClient *redis.Client, mongoClient *mongo.Client, environment *config.Environment) *Server {
	return &Server{database: database, redisClient: redisClient, mongoClient: mongoClient, environment: environment}
}

func (s *Server) ConfigureRoutes(applicationFactory *factory.ApplicationFactory) {
	server := echo.New()
	server.Use(middleware.Recover())
	server.Use(middleware.RequestID())
	server.Use(otelecho.Middleware("bom-pedido-api"))
	server.Use(middlewares.AuthenticateMiddleware(applicationFactory))
	server.HTTPErrorHandler = middlewares.HandleError

	api := server.Group("/api")
	shoppingCartRoutes := api.Group("/v1/shopping-cart", middlewares.LockByCustomerId(applicationFactory))
	shoppingCartRoutes.POST("/checkout", checkout_shopping_cart.Handle(applicationFactory))
	shoppingCartRoutes.PATCH("/items", add_item_to_shopping_cart.Handle(applicationFactory))

	api.POST("/v1/products", create_product.Handle(applicationFactory))
	api.GET("/v1/products", list_products.Handle(applicationFactory))
	api.POST("/v1/auth/google/customer", google_auth_customer.Handle(applicationFactory))
	api.GET("/v1/customers/me", get_customer.Handle(applicationFactory))

	orderRoutes := api.Group("/v1/orders/:id", middlewares.LockByParam("id", applicationFactory))
	orderRoutes.POST("/approve", approve.Handle(applicationFactory))
	orderRoutes.POST("/reject", reject.Handle(applicationFactory))
	orderRoutes.POST("/cancel", cancel.Handle(applicationFactory))
	orderRoutes.POST("/finish", finish.Handle(applicationFactory))
	orderRoutes.POST("/in-progress", mark_in_progress.Handle(applicationFactory))
	orderRoutes.POST("/delivering", mark_delivering.Handle(applicationFactory))
	orderRoutes.POST("/awaiting-withdraw", mark_awaiting_withdraw.Handle(applicationFactory))
	orderRoutes.POST("/awaiting-delivery", mark_awaiting_delivery.Handle(applicationFactory))
	server.GET("/api/health", health.Handle(s.database, s.redisClient, s.mongoClient))
	s.server = server
}

func (s *Server) StartTracer() {
	if os.Getenv("ENVIRONMENT") != "production" {
		return
	}
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(s.environment.OpenTelemetryEndpointExporter),
			otlptracehttp.WithHeaders(map[string]string{"content-type": "application/json"}),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		panic(err)
	}
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(
			exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
		),
		trace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String("bom-pedido-api"))),
	)
	otel.SetTracerProvider(tracerProvider)
	s.tracerProvider = tracerProvider
}

func (s *Server) Run(port string) {
	s.StartTracer()
	s.server.Logger.Fatal(s.server.Start(port))
}

func (s *Server) AwaitInterruptSignal() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop
}

func (s *Server) Shutdown() {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	if err := s.database.Close(); err != nil {
		slog.Error("Error on close database connection", "error", err)
	}
	if err := s.redisClient.Close(); err != nil {
		slog.Error("Error on close redis connection", "error", err)
	}
	if err := s.mongoClient.Disconnect(ctx); err != nil {
		slog.Error("Error on close mongo connection", "error", err)
	}
	if s.tracerProvider != nil {
		if err := s.tracerProvider.Shutdown(ctx); err != nil {
			slog.Error("Error on close trace provider connection", "error", err)
		}
	}
	slog.Info("Shutting down server...")
	if err := s.server.Shutdown(ctx); err != nil {
		slog.Error("Error on shutdown server", "error", err)
	}
}
