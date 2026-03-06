package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/profile/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/KimNattanan/go-chat-backend/pkg/rabbitmq"
	"github.com/go-playground/validator/v10"
)

// NewTranslationRoutes -.
func NewTranslationRoutes(routes map[string]rabbitmq.Handler, profileUseCase usecase.ProfileUseCase, l logger.Interface) {
	r := &V1{
		profileUsecase: profileUseCase,
		l:              l,
		v:              validator.New(validator.WithRequiredStructEnabled()),
	}

	{
		routes["createProfile"] = r.createProfile()
		routes["deleteProfile"] = r.deleteProfile()
	}
}
