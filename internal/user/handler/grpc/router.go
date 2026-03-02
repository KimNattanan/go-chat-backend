package rest

import (
	v1 "github.com/KimNattanan/go-chat-backend/internal/user/handler/grpc/v1"
	"github.com/KimNattanan/go-chat-backend/internal/user/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	pbgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// NewRouter -.
func NewRouter(app *pbgrpc.Server, profileUseCase usecase.ProfileUseCase, l logger.Interface) {
	{
		v1.NewUserRoutes(app, profileUseCase, l)
	}
	reflection.Register(app)
}
