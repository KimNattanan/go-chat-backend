package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/auth/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/go-playground/validator/v10"
)

// V1 -.
type V1 struct {
	authUseCase usecase.AuthUseCase
	l           logger.Interface
	v           *validator.Validate
	appEnv      string
}
