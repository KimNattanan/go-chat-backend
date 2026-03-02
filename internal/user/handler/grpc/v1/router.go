package v1

import (
	v1 "github.com/KimNattanan/go-chat-backend/internal/user/proto/v1"
	"github.com/KimNattanan/go-chat-backend/internal/user/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/go-playground/validator/v10"
	pbgrpc "google.golang.org/grpc"
)

// NewUserRoutes -.
func NewUserRoutes(app *pbgrpc.Server, profileUseCase usecase.ProfileUseCase, l logger.Interface) {
	r := &V1{
		profileUseCase: profileUseCase,
		l:              l,
		v:              validator.New(validator.WithRequiredStructEnabled()),
	}
	{
		v1.RegisterProfileServiceServer(app, r)
	}
}
