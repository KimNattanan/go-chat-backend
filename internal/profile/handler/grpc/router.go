package rest

import (
	v1 "github.com/KimNattanan/go-chat-backend/internal/profile/handler/grpc/v1"
	"github.com/KimNattanan/go-chat-backend/internal/profile/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	pbgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// NewRouter -.
func NewRouter(app *pbgrpc.Server, profileUseCase usecase.ProfileUseCase, l logger.Interface) {
	{
		v1.NewProfileRoutes(app, profileUseCase, l)
	}
	reflection.Register(app)
}
