package http

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/infra/config"
	"bom-pedido-api/internal/infra/http/admin"
	"bom-pedido-api/internal/infra/http/callback"
	"bom-pedido-api/internal/infra/http/category"
	"bom-pedido-api/internal/infra/http/customer"
	"bom-pedido-api/internal/infra/http/health"
	"bom-pedido-api/internal/infra/http/middlewares"
	"bom-pedido-api/internal/infra/http/order"
	"bom-pedido-api/internal/infra/http/products"
	"bom-pedido-api/internal/infra/http/shopping_cart"
	"bom-pedido-api/pkg/mongo"
	"context"
	"fmt"
	echoPrometheus "github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	server             *echo.Echo
	mongoDatabase      *mongo.Database
	redisClient        *redis.Client
	tracerProvider     *trace.TracerProvider
	environment        *config.Environment
	applicationFactory *factory.ApplicationFactory
}

func NewServer(
	redisClient *redis.Client,
	mongoDatabase *mongo.Database,
	environment *config.Environment,
	applicationFactory *factory.ApplicationFactory,
) *Server {
	return &Server{
		redisClient:        redisClient,
		mongoDatabase:      mongoDatabase,
		environment:        environment,
		applicationFactory: applicationFactory,
	}
}

func (s *Server) ConfigureRoutes() {
	server := echo.New()
	server.Use(middleware.Recover())
	server.Use(middleware.RequestID())
	server.Use(echoPrometheus.NewMiddleware("bom_pedido_api"))
	server.Use(otelecho.Middleware("bom-pedido-api"))
	server.Use(middlewares.AuthenticateMiddleware(s.applicationFactory))
	server.Use(middlewares.SetContextTenantId())
	server.HTTPErrorHandler = middlewares.HandleError

	server.GET("/metrics", echoPrometheus.NewHandler())

	api := server.Group("/api")
	admin.ConfigureRoutes(api, s.applicationFactory, s.environment)
	order.ConfigureRoutes(api, s.applicationFactory)
	products.ConfigureRoutes(api, s.applicationFactory)
	customer.ConfigureRoutes(api, s.applicationFactory)
	shopping_cart.ConfigureRoutes(api, s.applicationFactory)
	callback.ConfigureCallbackRoutes(api, s.applicationFactory)
	category.ConfigureRoutes(api, s.applicationFactory)

	server.GET("/api/health", health.Handle(s.redisClient, s.mongoDatabase))
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
		trace.WithResource(resource.NewWithAttributes(
			"https://opentelemetry.io/schemas/1.27.0",
			attribute.String("service.name", "bom-pedido-api"),
		)),
	)
	otel.SetTracerProvider(tracerProvider)
	s.tracerProvider = tracerProvider
}

func (s *Server) Run() {
	s.StartTracer()
	s.server.Logger.Fatal(s.server.Start(fmt.Sprintf(":%s", s.environment.Port)))
}

func (s *Server) AwaitInterruptSignal() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop
}

func (s *Server) Shutdown() {
	s.applicationFactory.Close()

	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	if err := s.redisClient.Close(); err != nil {
		slog.Error("Error on close redis connection", "error", err)
	}
	if err := s.mongoDatabase.Disconnect(ctx); err != nil {
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
