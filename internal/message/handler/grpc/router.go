package rest

import (
	v1 "github.com/KimNattanan/go-chat-backend/internal/message/handler/grpc/v1"
	"github.com/KimNattanan/go-chat-backend/internal/message/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	pbgrpc "google.golang.org/grpc"
)

// NewRouter -.
func NewRouter(app *pbgrpc.Server, messageUseCase usecase.MessageUseCase, l logger.Interface) {
	v1.NewMessageRoutes(app, messageUseCase, l)
}
