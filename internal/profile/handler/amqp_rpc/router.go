package v1

import (
	v1 "github.com/KimNattanan/go-chat-backend/internal/profile/handler/amqp_rpc/v1"
	"github.com/KimNattanan/go-chat-backend/internal/profile/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/KimNattanan/go-chat-backend/pkg/rabbitmq"
)

// NewRouter -.
func NewRouter(routes map[string]rabbitmq.Handler, profileUseCase usecase.ProfileUseCase, l logger.Interface) {
	v1.NewTranslationRoutes(routes, profileUseCase, l)
}
