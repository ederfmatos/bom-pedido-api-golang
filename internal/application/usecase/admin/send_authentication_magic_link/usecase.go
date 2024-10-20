package send_authentication_magic_link

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/application/token"
	"context"
	"fmt"
)

const (
	subject  = "Login no Bom Pedido"
	template = "sign-in-admin"
)

type (
	UseCase struct {
		adminRepository    repository.AdminRepository
		merchantRepository repository.MerchantRepository
		eventEmitter       event.Emitter
		tokenManager       token.Manager
		baseUrl            string
	}

	Input struct {
		Email string
	}
)

func New(baseUrl string, factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{
		adminRepository:    factory.AdminRepository,
		merchantRepository: factory.MerchantRepository,
		eventEmitter:       factory.EventEmitter,
		tokenManager:       factory.TokenManager,
		baseUrl:            baseUrl,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input Input) error {
	admin, err := uc.adminRepository.FindByEmail(ctx, input.Email)
	if err != nil {
		return err
	}
	if admin == nil {
		return nil
	}
	isActive, err := uc.merchantRepository.IsActive(ctx, admin.MerchantId)
	if err != nil || !isActive {
		return nil
	}
	tokenData := token.Data{
		Type: "MAGIC_LINK",
		Id:   admin.Id,
	}
	magicLinkToken, err := uc.tokenManager.Encrypt(ctx, tokenData)
	if err != nil {
		return nil
	}
	sendEmailEvent := event.NewSendEmailEvent(admin.Email, subject, map[string]string{
		"name":     admin.Name,
		"template": template,
		"url":      fmt.Sprintf("%s/%s", uc.baseUrl, magicLinkToken),
	})
	return uc.eventEmitter.Emit(ctx, sendEmailEvent)
}
