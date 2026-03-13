package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/profile/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/KimNattanan/go-chat-backend/pkg/rabbitmq"
	"github.com/go-playground/validator/v10"
)

// NewProfileRoutes -.
func NewProfileRoutes(profileUseCase usecase.ProfileUseCase, l logger.Interface) map[string]rabbitmq.Handler {
	r := &V1{
		profileUsecase: profileUseCase,
		l:              l,
		v:              validator.New(validator.WithRequiredStructEnabled()),
	}
	routes := make(map[string]rabbitmq.Handler)
	routes["user.created"] = r.createProfile
	routes["user.deleted"] = r.deleteProfile
	return routes
}
