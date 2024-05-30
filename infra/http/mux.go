package http

import (
	"bom-pedido-api/presentation/rest"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type MuxHttpRequest struct {
	body map[string]any
	*http.Request
}

func NewMuxHttpRequest(request *http.Request) rest.Request {
	var body map[string]any
	decoder := json.NewDecoder(request.Body)
	_ = decoder.Decode(&body)
	return &MuxHttpRequest{
		body:    body,
		Request: request,
	}
}

func (request *MuxHttpRequest) Context() context.Context {
	return request.Context()
}

func (request *MuxHttpRequest) BodyFieldString(name string) string {
	value := request.body[name]
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%s", value)
}

func (request *MuxHttpRequest) BodyFieldFloat(name string) float64 {
	value := request.body[name]
	if value == nil {
		return 0.0
	}
	floatValue, _ := strconv.ParseFloat(fmt.Sprintf("%f", value), 64)
	return floatValue
}
