package gateway

import (
	"bom-pedido-api/application/gateway"
	"bom-pedido-api/infra/json"
	"context"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

type DefaultGoogleGateway struct {
	baseUrl string
}

func NewDefaultGoogleGateway(baseUrl string) gateway.GoogleGateway {
	return &DefaultGoogleGateway{baseUrl: baseUrl}
}

func (googleGateway *DefaultGoogleGateway) GetUserByToken(ctx context.Context, token string) (*gateway.GoogleUser, error) {
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	request, err := http.NewRequestWithContext(ctx, "GET", googleGateway.baseUrl+"?access_token="+token, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	var googleUser gateway.GoogleUser
	err = json.Decode(ctx, response.Body, &googleUser)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return &googleUser, nil
}
