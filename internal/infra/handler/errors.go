package handler

import (
	"bom-pedido-api/internal/application/usecase/category"
	domainErrors "bom-pedido-api/internal/domain/errors"
	"bom-pedido-api/pkg/http"
	"errors"
)

var ErrInvalidCurrentPage = errors.New("invalid current page")
var ErrInvalidPageSize = errors.New("invalid page size")

func ErrorHandlerMiddleware() http.Middleware {
	var mappedErrors = map[error]http.MappedError{
		domainErrors.OrderNotFoundError: {
			Status:   http.StatusNotFound,
			Response: http.NewErrorResponse("Pedido não encontrado"),
		},
		domainErrors.CustomerNotFoundError: {
			Status:   http.StatusNotFound,
			Response: http.NewErrorResponse("Cliente não encontrado"),
		},
		domainErrors.ProductNotFoundError: {
			Status:   http.StatusNotFound,
			Response: http.NewErrorResponse("Produto não encontrado"),
		},
		domainErrors.ProductPriceIsRequiredError: {
			Status:   http.StatusBadRequest,
			Response: http.NewErrorResponse("O preço do produto é obrigatório"),
		},
		domainErrors.ProductPriceShouldPositiveError: {
			Status:   http.StatusBadRequest,
			Response: http.NewErrorResponse("O preço do produto deve ser maior que zero"),
		},
		domainErrors.ProductNameIsRequiredError: {
			Status:   http.StatusBadRequest,
			Response: http.NewErrorResponse("O nome do produto é obrigatório"),
		},
		domainErrors.ProductUnAvailableError: {
			Status:   http.StatusConflict,
			Response: http.NewErrorResponse("Produto indisponível"),
		},
		domainErrors.ProductWithSameNameError: {
			Status:   http.StatusConflict,
			Response: http.NewErrorResponse("Já existe um produto com esse nome"),
		},
		domainErrors.ProductCategoryNotFoundError: {
			Status:   http.StatusNotFound,
			Response: http.NewErrorResponse("Categoria de produto não encontrada"),
		},
		domainErrors.ProductInvalidProductStatusError: {
			Status:   http.StatusBadRequest,
			Response: http.NewErrorResponse("Status de produto inválido"),
		},
		domainErrors.OrderDeliveryModeIsWithdrawError: {
			Status:   http.StatusConflict,
			Response: http.NewErrorResponse("Esse pedido está para ser retirado, a entrega à domicilio não é permitida"),
		},
		domainErrors.OrderDeliveryModeIsDeliveryError: {
			Status:   http.StatusConflict,
			Response: http.NewErrorResponse("Esse pedido está para ser entregue a domicilop, a retirada não é permitida"),
		},
		domainErrors.ShoppingCartEmptyError: {
			Status:   http.StatusUnprocessableEntity,
			Response: http.NewErrorResponse("Seu carrinho está vazio"),
		},
		domainErrors.CardTokenIsRequiredError: {
			Status:   http.StatusBadRequest,
			Response: http.NewErrorResponse("O token do cartão é obrigatório"),
		},
		domainErrors.PaybackShouldBePositiveError: {
			Status:   http.StatusBadRequest,
			Response: http.NewErrorResponse("O valor do troco deve ser maior que zero"),
		},
		category.CategoryWithSameNameError: {
			Status:   http.StatusConflict,
			Response: http.NewErrorResponse("Já existe uma categoria criada com esse nome"),
		},
		ErrInvalidCurrentPage: {
			Status:   http.StatusBadRequest,
			Response: http.NewErrorResponse("Valor inválido para página atual"),
		},
		ErrInvalidPageSize: {
			Status:   http.StatusBadRequest,
			Response: http.NewErrorResponse("Valor inválido para tamanho da página"),
		},
	}
	return http.ErrorHandlerHTTPMiddleware(mappedErrors)
}
