package rest

import (
	v1 "github.com/KimNattanan/go-chat-backend/internal/auth/handler/grpc/v1"
	"github.com/KimNattanan/go-chat-backend/internal/auth/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	pbgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// NewRouter -.
func NewRouter(app *pbgrpc.Server, authUseCase usecase.AuthUseCase, l logger.Interface) {
	v1.NewAuthRoutes(app, authUseCase, l)
	reflection.Register(app)
}
