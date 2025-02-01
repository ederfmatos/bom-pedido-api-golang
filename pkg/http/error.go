package http

import (
	"bom-pedido-api/pkg/log"
	"errors"
)

type (
	ErrorResponse struct {
		Error string `json:"error"`
	}

	MappedError struct {
		Status   Status
		Response interface{}
		Handler  func(err error) ErrorResponse
	}
)

func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{Error: message}
}

func ErrorHandlerHTTPMiddleware(mappedErrors map[error]MappedError) Middleware {
	return func(handler Handler) Handler {
		return func(request Request, response Response) error {
			responseError := handler(request, response)
			if responseError == nil {
				return nil
			}

			log.Error("Ocorreu um erro na requisição", responseError, "request", request.String())

			mappedError, found := mappedErrors[responseError]
			if !found {
				for err, mappedErr := range mappedErrors {
					if errors.Is(responseError, err) {
						mappedError = mappedErr
						break
					}
				}
				if mappedError.Status == 0 {
					return responseError
				}
			}

			var responseBody = mappedError.Response
			if mappedError.Handler != nil {
				responseBody = mappedError.Handler(responseError)
			}

			response.SetStatus(mappedError.Status)
			return response.SetBody(responseBody)
		}
	}
}
