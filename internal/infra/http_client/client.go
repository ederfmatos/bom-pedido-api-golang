package http_client

import (
	"context"
)

type (
	HTTPClient interface {
		Post(path string, values ...string) HTTPClientBuilder
		Get(path string, values ...string) HTTPClientBuilder
	}

	HTTPClientBuilder interface {
		Body(body interface{}) HTTPClientBuilder
		Header(key string, value string) HTTPClientBuilder
		Execute(ctx context.Context) (HTTPResponse, error)
	}

	HTTPResponse interface {
		Close()
		IsError() bool
		ParseBody(value interface{}) error
		ParseError(value error) error
		GetErrorMessage() string
	}
)
