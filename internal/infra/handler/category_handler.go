package handler

import (
	usecase "bom-pedido-api/internal/application/usecase/category"
	"bom-pedido-api/pkg/http"
	"fmt"
)

type (
	CategoryHandler struct {
		createCategoryUseCase *usecase.CreateCategoryUseCase
	}

	createCategoryBody struct {
		Name        string `body:"name" json:"name,omitempty"`
		Description string `body:"description" json:"description,omitempty"`
	}
)

func NewCategoryHandler(createCategoryUseCase *usecase.CreateCategoryUseCase) *CategoryHandler {
	return &CategoryHandler{createCategoryUseCase: createCategoryUseCase}
}

func (h CategoryHandler) CreateCategory(request http.Request, response http.Response) error {
	var body createCategoryBody
	if err := request.Bind(&body); err != nil {
		return err
	}

	input := usecase.CreateCategoryInput{
		Name:        body.Name,
		Description: body.Description,
		TenantId:    request.TenantID(),
	}
	output, err := h.createCategoryUseCase.Execute(request.Context(), input)
	if err != nil {
		return fmt.Errorf("create category: %w", err)
	}

	return response.Created(output)
}
