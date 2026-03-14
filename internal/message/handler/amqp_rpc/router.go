package amqp_rpc

import (
	v1 "github.com/KimNattanan/go-chat-backend/internal/message/handler/amqp_rpc/v1"
	"github.com/KimNattanan/go-chat-backend/internal/message/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/KimNattanan/go-chat-backend/pkg/rabbitmq"
)

// NewRouter -.
func NewRouter(messageUseCase usecase.MessageUseCase, l logger.Interface) map[string]rabbitmq.Handler {
	return v1.NewMessageRoutes(messageUseCase, l)
}
