package net

import (
	"bom-pedido-api/internal/infra/json"
	"bom-pedido-api/pkg/http"
	"context"
	"fmt"
	net "net/http"
	"strconv"
	"strings"
)

const (
	_tenantID = "TENANT_ID"
	_userID   = "USER_ID"
)

type Request struct {
	request *net.Request
}

func NewRequest(request *net.Request) http.Request {
	return &Request{
		request: request,
	}
}

func (r Request) QueryParam() http.Query {
	return r
}

func (r Request) Context() context.Context {
	return r.request.Context()
}

func (r Request) WithContext(ctx context.Context) http.Request {
	return NewRequest(r.request.WithContext(ctx))
}

func (r Request) TenantID() string {
	tenantID := r.Context().Value(_tenantID)
	if tenantID != nil {
		return tenantID.(string)
	}

	return strings.Split(r.request.Host, ":")[0]
}

func (r Request) PathParam(name string) string {
	return r.request.PathValue(name)
}

func (r Request) AuthenticatedUser() string {
	if userID := r.Context().Value(_userID); userID != nil {
		return userID.(string)
	}
	return ""
}

func (r Request) Bind(target interface{}) error {
	if err := json.Decode(r.Context(), r.request.Body, target); err != nil {
		return fmt.Errorf("decode body: %w", err)
	}

	return nil
}

func (r Request) GetHeader(name http.Header) string {
	return r.request.Header.Get(string(name))
}

func (r Request) SetHeader(name http.Header, value string) {
	r.request.Header.Set(string(name), value)
}

func (r Request) Int(name string) (int, error) {
	value := r.request.URL.Query().Get(name)
	return strconv.Atoi(value)
}

func (r Request) String() string {
	return fmt.Sprintf("%s %s", r.request.Method, r.request.URL.Path)
}
