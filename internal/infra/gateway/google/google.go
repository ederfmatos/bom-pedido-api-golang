package google

import (
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/internal/infra/http_client"
	"context"
	"errors"
)

type DefaultGoogleGateway struct {
	httpClient http_client.HTTPClient
}

func NewDefaultGoogleGateway(client http_client.HTTPClient) gateway.GoogleGateway {
	return &DefaultGoogleGateway{httpClient: client}
}

func (googleGateway *DefaultGoogleGateway) GetUserByToken(ctx context.Context, token string) (*gateway.GoogleUser, error) {
	httpResponse, err := googleGateway.httpClient.Get("?access_token=", token).Execute(ctx)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Close()
	if httpResponse.IsError() {
		return nil, errors.New(httpResponse.GetErrorMessage())
	}
	var googleUser gateway.GoogleUser
	if err = httpResponse.ParseBody(&googleUser); err != nil {
		return nil, err
	}
	return &googleUser, nil
}
