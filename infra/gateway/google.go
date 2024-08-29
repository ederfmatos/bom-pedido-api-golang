package gateway

import (
	"bom-pedido-api/application/gateway"
	"bom-pedido-api/infra/http_client"
	"context"
	"errors"
)

type DefaultGoogleGateway struct {
	httpClient http_client.HttpClient
}

func NewDefaultGoogleGateway(client http_client.HttpClient) gateway.GoogleGateway {
	return &DefaultGoogleGateway{httpClient: client}
}

func (googleGateway *DefaultGoogleGateway) GetUserByToken(ctx context.Context, token string) (*gateway.GoogleUser, error) {
	httpResponse, err := googleGateway.httpClient.Get("?access_token=", token).Execute(ctx)
	if err != nil {
		return nil, err
	}
	if httpResponse.IsError() {
		return nil, errors.New(httpResponse.GetErrorMessage())
	}
	var googleUser gateway.GoogleUser
	if err = httpResponse.ParseBody(&googleUser); err != nil {
		return nil, err
	}
	return &googleUser, nil
}
