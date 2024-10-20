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
	"context"
	"database/sql"
	echoPrometheus "github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	echoSwagger "github.com/swaggo/echo-swagger"
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
	server.Use(middlewares.RedocDocumentation())
	server.Use(echoPrometheus.NewMiddleware("bom_pedido_api"))
	server.Use(otelecho.Middleware("bom-pedido-api"))
	server.Use(middlewares.AuthenticateMiddleware(applicationFactory))
	server.Use(middlewares.SetContextTenantId())
	server.HTTPErrorHandler = middlewares.HandleError

	server.GET("/metrics", echoPrometheus.NewHandler())
	server.GET("/swagger.json", func(c echo.Context) error {
		return c.File(".docs/openapi.json")
	})
	server.GET("/swagger/*", echoSwagger.EchoWrapHandler(func(c *echoSwagger.Config) {
		c.URLs = []string{".docs/openapi.json"}
	}))

	api := server.Group("/api")
	admin.ConfigureRoutes(api, applicationFactory, s.environment)
	order.ConfigureRoutes(api, applicationFactory)
	products.ConfigureRoutes(api, applicationFactory)
	customer.ConfigureRoutes(api, applicationFactory)
	shopping_cart.ConfigureRoutes(api, applicationFactory)
	callback.ConfigureCallbackRoutes(api, applicationFactory)
	category.ConfigureRoutes(api, applicationFactory)

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
