package handler

import (
	usecase "bom-pedido-api/internal/application/usecase/customer"
	"bom-pedido-api/pkg/http"
	"fmt"
)

type (
	CustomerHandler struct {
		getCustomerUseCase          *usecase.GetCustomerUseCase
		authenticateCustomerUseCase *usecase.GoogleAuthenticateCustomerUseCase
	}

	authenticateCustomerBody struct {
		Token string `json:"token"`
	}
)

func NewCustomerHandler(
	getCustomerUseCase *usecase.GetCustomerUseCase,
	authenticateCustomerUseCase *usecase.GoogleAuthenticateCustomerUseCase,
) *CustomerHandler {
	return &CustomerHandler{
		getCustomerUseCase:          getCustomerUseCase,
		authenticateCustomerUseCase: authenticateCustomerUseCase,
	}
}

func (h CustomerHandler) GetCustomer(request http.Request, response http.Response) error {
	input := usecase.GetCustomerInput{
		Id: request.AuthenticatedUser(),
	}
	output, err := h.getCustomerUseCase.Execute(request.Context(), input)
	if err != nil {
		return fmt.Errorf("get customer: %w", err)
	}
	return response.OK(output)
}

func (h CustomerHandler) AuthenticateCustomer(request http.Request, response http.Response) error {
	var body authenticateCustomerBody
	if err := request.Bind(&body); err != nil {
		return err
	}

	input := usecase.GoogleAuthenticateCustomerInput{
		Token:    body.Token,
		TenantId: request.TenantID(),
	}
	output, err := h.authenticateCustomerUseCase.Execute(request.Context(), input)
	if err != nil {
		return fmt.Errorf("authenticate customer: %w", err)
	}

	return response.OK(output)
}
