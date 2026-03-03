package v1

import (
	v1 "github.com/KimNattanan/go-chat-backend/internal/profile/proto/v1"
	"github.com/KimNattanan/go-chat-backend/internal/profile/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/go-playground/validator/v10"
	pbgrpc "google.golang.org/grpc"
)

// NewProfileRoutes -.
func NewProfileRoutes(app *pbgrpc.Server, profileUseCase usecase.ProfileUseCase, l logger.Interface) {
	r := &V1{
		profileUseCase: profileUseCase,
		l:              l,
		v:              validator.New(validator.WithRequiredStructEnabled()),
	}
	{
		v1.RegisterProfileServiceServer(app, r)
	}
}
