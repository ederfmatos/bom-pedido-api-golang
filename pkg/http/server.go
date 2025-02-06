package http

import (
	"bom-pedido-api/pkg/telemetry"
	"context"
)

type (
	Query interface {
		Int(name string) (int, error)
	}

	Request interface {
		QueryParam() Query
		Context() context.Context
		WithContext(ctx context.Context) Request
		TenantID() string
		PathParam(name string) string
		AuthenticatedUser() string
		Bind(target interface{}) error
		String() string
		GetHeader(name Header) string
		SetHeader(name Header, value string)
		Method() string
		URL() string
	}

	Response interface {
		Bytes() []byte
		Status() Status
		Headers() map[string]string
		SetStatus(status Status)
		SetBody(body interface{}) error
		OK(body interface{}) error
		Created(body interface{}) error
		NoContent() error
	}

	Handler func(request Request, response Response) error

	Middleware func(handler Handler) Handler

	Server interface {
		Post(path string, handler Handler, middlewares ...Middleware)
		Get(path string, handler Handler, middlewares ...Middleware)
		Patch(path string, handler Handler, middlewares ...Middleware)
		Put(path string, handler Handler, middlewares ...Middleware)
		Delete(path string, handler Handler, middlewares ...Middleware)
		AddMiddleware(middleware Middleware)
		UseTracerProvider(provider telemetry.TracerProvider)
		AwaitInterruptSignal()
		Shutdown()
		Run()
	}
)
