package amqp_rpc

import (
	v1 "github.com/KimNattanan/go-chat-backend/internal/chat/handler/amqp_rpc/v1"
	"github.com/KimNattanan/go-chat-backend/internal/chat/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/KimNattanan/go-chat-backend/pkg/rabbitmq"
)

// NewRouter -.
func NewRouter(roomUseCase usecase.RoomUseCase, membershipUseCase usecase.MembershipUseCase, l logger.Interface) map[string]rabbitmq.Handler {
	return v1.NewChatRoutes(roomUseCase, membershipUseCase, l)
}
