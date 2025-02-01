package handler

import (
	usecase "bom-pedido-api/internal/application/usecase/admin"
	"bom-pedido-api/pkg/http"
	"fmt"
)

type (
	AdminHandler struct {
		sendAuthenticationMagicLinkUseCase *usecase.SendAuthenticationMagicLinkUseCase
	}

	sendAuthenticationLinkBody struct {
		Email string `json:"email"`
	}
)

func NewAdminHandler(sendAuthenticationMagicLinkUseCase *usecase.SendAuthenticationMagicLinkUseCase) *AdminHandler {
	return &AdminHandler{sendAuthenticationMagicLinkUseCase: sendAuthenticationMagicLinkUseCase}
}

func (h AdminHandler) SendAuthenticationLink(request http.Request, response http.Response) error {
	var body sendAuthenticationLinkBody
	if err := request.Bind(&body); err != nil {
		return err
	}

	input := usecase.SendAuthenticationMagicLinkInput{Email: body.Email}
	if err := h.sendAuthenticationMagicLinkUseCase.Execute(request.Context(), input); err != nil {
		return fmt.Errorf("send authentication link: %w", err)
	}

	return response.NoContent()
}
