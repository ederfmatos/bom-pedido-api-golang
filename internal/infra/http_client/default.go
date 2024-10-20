package http_client

import (
	"bom-pedido-api/internal/infra/json"
	"bytes"
	"context"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"io"
	"net/http"
	"strings"
)

type (
	defaultHttpClient struct {
		baseUrl string
		client  http.Client
	}
	defaultHttpClientBuilder struct {
		method  string
		url     string
		headers map[string]string
		body    *interface{}
		client  http.Client
	}
	defaultHttpResponse struct {
		*http.Response
		ctx context.Context
	}
)

func NewDefaultHttpClient(baseUrl string) HTTPClient {
	return &defaultHttpClient{
		baseUrl: baseUrl,
		client:  http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)},
	}
}

func newDefaultHttpClientBuilder(client http.Client, method, url string, values ...string) HTTPClientBuilder {
	return &defaultHttpClientBuilder{
		method:  method,
		url:     url + strings.Join(values, ""),
		headers: make(map[string]string),
		body:    nil,
		client:  client,
	}
}

func (client *defaultHttpClient) Post(path string, values ...string) HTTPClientBuilder {
	return newDefaultHttpClientBuilder(client.client, "POST", client.baseUrl+path, values...)
}

func (client *defaultHttpClient) Get(path string, values ...string) HTTPClientBuilder {
	return newDefaultHttpClientBuilder(client.client, "GET", client.baseUrl+path, values...)
}

func (builder *defaultHttpClientBuilder) Body(body interface{}) HTTPClientBuilder {
	builder.body = &body
	return builder
}

func (builder *defaultHttpClientBuilder) Header(key string, value string) HTTPClientBuilder {
	builder.headers[key] = value
	return builder
}

func (builder *defaultHttpClientBuilder) Execute(ctx context.Context) (HTTPResponse, error) {
	var body io.Reader
	if builder.body != nil {
		paymentBytes, err := json.Marshal(ctx, builder.body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(paymentBytes)
	} else {
		body = nil
	}
	request, err := http.NewRequestWithContext(ctx, builder.method, builder.url, body)
	if err != nil {
		return nil, err
	}
	for key, value := range builder.headers {
		request.Header.Add(key, value)
	}
	response, err := builder.client.Do(request)
	if err != nil {
		return nil, err
	}
	return &defaultHttpResponse{Response: response, ctx: ctx}, nil
}

func (r *defaultHttpResponse) IsError() bool {
	return r.StatusCode >= 400
}

func (r *defaultHttpResponse) Close() {
	_ = r.Body.Close()
}

func (r *defaultHttpResponse) ParseBody(value interface{}) error {
	return json.Decode(r.ctx, r.Body, value)
}

func (r *defaultHttpResponse) ParseError(value error) error {
	if err := json.Decode(r.ctx, r.Body, value); err != nil {
		return err
	}
	return value
}

func (r *defaultHttpResponse) GetErrorMessage() string {
	var mapResponse map[string]interface{}
	_ = json.Decode(r.ctx, r.Body, &mapResponse)
	var errorMessage interface{}
	if value := mapResponse["error"]; value != nil {
		errorMessage = value
	} else if value = mapResponse["message"]; value != nil {
		errorMessage = value
	} else {
		return r.Status
	}
	if message, ok := errorMessage.(string); ok {
		return message
	}
	return r.Status
}
