package rest

import (
	"context"
)

type Request interface {
	BodyFieldString(name string) string
	BodyFieldFloat(name string) float64
	Context() context.Context
}

type ResponseWriter interface {
	StatusOk(body interface{}) error
	StatusNoContent() ResponseWriter
	Status(status int) ResponseWriter
	Body(body interface{}) error
	Header(string, string) ResponseWriter
	HandleError(err error) error
}

type RequestHandler func(request Request, responseWriter ResponseWriter) error
