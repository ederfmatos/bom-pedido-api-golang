package gateway

import (
	"bom-pedido-api/application/gateway"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type DefaultGoogleGateway struct {
	baseUrl string
}

func NewDefaultGoogleGateway(baseUrl string) gateway.GoogleGateway {
	return &DefaultGoogleGateway{baseUrl: baseUrl}
}

func (googleGateway *DefaultGoogleGateway) GetUserByToken(ctx context.Context, token string) (*gateway.GoogleUser, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", googleGateway.baseUrl+"?access_token="+token, nil)
	if err != nil {
		return nil, err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var googleUser gateway.GoogleUser
	err = json.Unmarshal(responseData, &googleUser)
	if err != nil {
		return nil, err
	}
	return &googleUser, nil
}
