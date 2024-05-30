package handler

import (
	"bom-pedido-api/application/usecase/auth"
	"bom-pedido-api/infra/registry"
	"bom-pedido-api/presentation/rest"
)

func GoogleAuthCustomerHandler(request rest.Request, responseWriter rest.ResponseWriter) error {
	input := auth.Input{
		Token:   request.BodyFieldString("token"),
		Context: request.Context(),
	}
	useCase := registry.GetDependency[auth.GoogleAuthenticateCustomerUseCase]("GoogleAuthenticateCustomerUseCase")
	output, err := useCase.Execute(input)
	if err != nil {
		return responseWriter.HandleError(err)
	}
	return responseWriter.StatusOk(output)
}
