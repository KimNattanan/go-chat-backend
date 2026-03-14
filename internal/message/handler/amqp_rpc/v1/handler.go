package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/message/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/go-playground/validator/v10"
)

// V1 -.
type V1 struct {
	messageUseCase usecase.MessageUseCase
	l logger.Interface
	v *validator.Validate
}
