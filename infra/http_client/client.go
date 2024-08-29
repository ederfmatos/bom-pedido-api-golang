package http_client

import (
	"context"
)

type (
	HttpClient interface {
		Post(path string, values ...string) HttpClientBuilder
		Get(path string, values ...string) HttpClientBuilder
	}

	HttpClientBuilder interface {
		Body(body interface{}) HttpClientBuilder
		Header(key string, value string) HttpClientBuilder
		Execute(ctx context.Context) (HttpResponse, error)
	}

	HttpResponse interface {
		IsError() bool
		ParseBody(value interface{}) error
		GetErrorMessage() string
	}
)
