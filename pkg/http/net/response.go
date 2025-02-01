package net

import (
	"bom-pedido-api/internal/infra/json"
	"bom-pedido-api/pkg/http"
	"context"
	"fmt"
)

type Response struct {
	responseBytes []byte
	status        http.Status
	headers       map[string]string
}

func NewResponse() *Response {
	return &Response{
		headers:       make(map[string]string),
		status:        http.StatusOK,
		responseBytes: make([]byte, 0),
	}
}

func (r *Response) Bytes() []byte {
	return r.responseBytes
}

func (r *Response) Status() http.Status {
	return r.status
}

func (r *Response) Headers() map[string]string {
	if r.headers == nil {
		r.headers = make(map[string]string)
	}
	return r.headers
}

func (r *Response) SetBody(body interface{}) error {
	responseBytes, err := json.Marshal(context.Background(), body)
	if err != nil {
		return fmt.Errorf("marshal response: %w", err)
	}
	r.responseBytes = responseBytes
	return nil
}

func (r *Response) SetStatus(status http.Status) {
	r.status = status
}

func (r *Response) OK(body interface{}) error {
	r.status = http.StatusOK
	return r.SetBody(body)
}

func (r *Response) Created(body interface{}) error {
	r.status = http.StatusCreated
	return r.SetBody(body)
}

func (r *Response) NoContent() error {
	r.status = http.StatusNoContent
	return nil
}
