package handler

import (
	"bom-pedido-api/internal/application/projection"
	"bom-pedido-api/internal/application/query"
	usecase "bom-pedido-api/internal/application/usecase/product"
	"bom-pedido-api/pkg/http"
	"fmt"
)

type (
	ProductHandler struct {
		createProductUseCase *usecase.CreateProductUseCase
		productQuery         query.ProductQuery
	}

	createProductBody struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		CategoryId  string  `json:"categoryId"`
		Price       float64 `json:"price"`
	}
)

func NewProductHandler(createProductUseCase *usecase.CreateProductUseCase, productQuery query.ProductQuery) *ProductHandler {
	return &ProductHandler{createProductUseCase: createProductUseCase, productQuery: productQuery}
}

func (h ProductHandler) CreateProduct(request http.Request, response http.Response) error {
	var body createProductBody
	if err := request.Bind(&body); err != nil {
		return err
	}

	input := usecase.CreateProductInput{
		Name:        body.Name,
		Description: body.Description,
		CategoryId:  body.CategoryId,
		Price:       body.Price,
		TenantId:    request.TenantID(),
	}
	output, err := h.createProductUseCase.Execute(request.Context(), input)
	if err != nil {
		return fmt.Errorf("create product: %w", err)
	}

	return response.Created(output)
}

func (h ProductHandler) ListProducts(request http.Request, response http.Response) error {
	currentPage, err := request.QueryParam().Int("currentPage")
	if err != nil {
		return ErrInvalidCurrentPage
	}

	pageSize, err := request.QueryParam().Int("pageSize")
	if err != nil {
		return ErrInvalidPageSize
	}

	filter := projection.ProductListFilter{
		CurrentPage: int64(currentPage),
		PageSize:    int64(pageSize),
		TenantId:    request.TenantID(),
	}
	products, err := h.productQuery.List(request.Context(), filter)
	if err != nil {
		return fmt.Errorf("list products: %w", err)
	}

	return response.OK(products)
}
