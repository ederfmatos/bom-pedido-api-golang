package http

import (
	"bom-pedido-api/internal/application/token"
	"bom-pedido-api/pkg/log"
	"bom-pedido-api/pkg/telemetry"
	"context"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

const (
	_tenantID  = "TENANT_ID"
	_tokenType = "TOKEN_TYPE"
	_userID    = "USER_ID"

	MessageYouNeedToBeAuthenticated = "Você precisa estar autenticado para realizar essa operação"
)

type (
	locker interface {
		Lock(ctx context.Context, key ...string) (string, error)
		Release(ctx context.Context, key string)
	}

	tokenManager interface {
		Decrypt(ctx context.Context, token string) (*token.Data, error)
	}

	Middlewares struct {
		locker       locker
		tokenManager tokenManager
	}
)

func NewMiddlewares(locker locker, tokenManager tokenManager) *Middlewares {
	return &Middlewares{
		locker:       locker,
		tokenManager: tokenManager,
	}
}

func (m Middlewares) LockByRequestPath(pathName string) Middleware {
	return func(handler Handler) Handler {
		return func(request Request, response Response) error {
			param := request.PathParam(pathName)

			lockKey, err := m.locker.Lock(request.Context(), param)
			if err != nil {
				return fmt.Errorf("lock request by path: %v", err)
			}

			defer m.locker.Release(context.Background(), lockKey)
			return handler(request, response)
		}
	}
}

func (m Middlewares) AuthenticateMiddleware() Middleware {
	return func(handler Handler) Handler {
		return func(request Request, response Response) error {
			authorization := request.GetHeader(HeaderAuthorization)
			if authorization == "" {
				return handler(request, response)
			}

			authorization = strings.ReplaceAll(authorization, "Bearer ", "")

			ctx := request.Context()
			tokenData, err := m.tokenManager.Decrypt(ctx, authorization)
			if err != nil {
				response.SetStatus(StatusUnAuthorized)
				return fmt.Errorf("decrypt token: %w", err)
			}

			ctx = context.WithValue(ctx, _userID, tokenData.Id)
			ctx = context.WithValue(ctx, _tokenType, tokenData.Type)
			ctx = context.WithValue(ctx, _tenantID, tokenData.TenantId)

			return handler(request.WithContext(ctx), response)
		}
	}
}

func (m Middlewares) OnlyAdminMiddleware(handler Handler) Handler {
	return func(request Request, response Response) error {
		user := request.AuthenticatedUser()
		if user == "" {
			response.SetStatus(StatusUnAuthorized)
			return response.SetBody(NewErrorResponse(MessageYouNeedToBeAuthenticated))
		}
		return handler(request, response)
	}
}

func (m Middlewares) TelemetryMiddleware(handler Handler) Handler {
	return func(request Request, response Response) error {
		return telemetry.StartSpanReturningError(request.Context(), request.String(), func(ctx context.Context) error {
			return handler(request.WithContext(ctx), response)
		}, "http.method", request.Method(), "http.url", request.URL(), "tenant.id", request.TenantID())
	}
}

func (m Middlewares) OnlyCustomerMiddleware(handler Handler) Handler {
	return func(request Request, response Response) error {
		user := request.AuthenticatedUser()
		if user == "" {
			response.SetStatus(StatusUnAuthorized)
			return response.SetBody(NewErrorResponse(MessageYouNeedToBeAuthenticated))
		}
		return handler(request, response)
	}
}

func (m Middlewares) RecoverMiddleware(handler Handler) Handler {
	return func(request Request, response Response) error {
		defer m.recoverRequest(request)
		return handler(request, response)
	}
}

func (m Middlewares) recoverRequest(request Request) {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}
		log.Error("panic recovered", err, "request", request.String())
	}
}

func (m Middlewares) RequestIDMiddleware(handler Handler) Handler {
	return func(request Request, response Response) error {
		requestID := request.GetHeader(HeaderRequestID)
		if requestID == "" {
			requestID = uuid.NewString()
		}
		request.SetHeader(HeaderRequestID, requestID)

		return handler(request, response)
	}
}
