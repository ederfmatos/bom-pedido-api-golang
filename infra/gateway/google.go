package gateway

import (
	"bom-pedido-api/application/gateway"
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

func (googleGateway *DefaultGoogleGateway) GetUserByToken(token string) (*gateway.GoogleUser, error) {
	response, err := http.Get(googleGateway.baseUrl + "?access_token=" + token)
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
