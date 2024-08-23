package middlewares

import "C"
import (
	"bom-pedido-api/infra/tenant"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"strings"
)

const tenantIdHeader = "X-Tenant-Id"

func SetContextTenantId() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()
			path := request.URL.Path
			if !strings.HasPrefix(path, "/api") {
				return next(c)
			}
			tenantId := request.Header.Get(tenantIdHeader)
			if tenantId == "" {
				return fmt.Errorf("%s header missing", tenantIdHeader)
			}
			ctx := request.Context()
			span := trace.SpanFromContext(ctx)
			span.SetAttributes(attribute.String("tenant.id", tenantId))
			c.SetRequest(request.WithContext(context.WithValue(ctx, tenant.Id, tenantId)))
			c.Set(tenant.Id, tenantId)
			return next(c)
		}
	}
}
