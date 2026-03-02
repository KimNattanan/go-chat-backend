package v1

import (
	v1 "github.com/KimNattanan/go-chat-backend/internal/user/proto/v1"
	"github.com/KimNattanan/go-chat-backend/internal/user/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/go-playground/validator/v10"
)

// V1 -.
type V1 struct {
	v1.ProfileServiceServer

	profileUseCase usecase.ProfileUseCase
	l              logger.Interface
	v              *validator.Validate
}
