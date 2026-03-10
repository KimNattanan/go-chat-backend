package v1

import (
	v1 "github.com/KimNattanan/go-chat-backend/internal/message/proto/v1"
	"github.com/KimNattanan/go-chat-backend/internal/message/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/go-playground/validator/v10"
	pbgrpc "google.golang.org/grpc"
)

// NewMessageRoutes -.
func NewMessageRoutes(app *pbgrpc.Server, messageUseCase usecase.MessageUseCase, l logger.Interface) {
	r := &V1{
		messageUseCase: messageUseCase,
		l:              l,
		v:              validator.New(validator.WithRequiredStructEnabled()),
	}
	{
		v1.RegisterMessageServiceServer(app, r)
	}
}
