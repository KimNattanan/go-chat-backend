package v1

import (
	v1 "github.com/KimNattanan/go-chat-backend/internal/auth/proto/v1"
	"github.com/KimNattanan/go-chat-backend/internal/auth/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/go-playground/validator/v10"
	pbgrpc "google.golang.org/grpc"
)

// NewAuthRoutes -.
func NewAuthRoutes(app *pbgrpc.Server, authUseCase usecase.AuthUseCase, l logger.Interface) {
	r := &V1{
		authUseCase:    authUseCase,
		l:              l,
		v:              validator.New(validator.WithRequiredStructEnabled()),
	}
	{
		v1.RegisterAuthServiceServer(app, r)
	}
}
