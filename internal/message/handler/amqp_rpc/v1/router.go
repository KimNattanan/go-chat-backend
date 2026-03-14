package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/message/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/KimNattanan/go-chat-backend/pkg/rabbitmq"
	"github.com/go-playground/validator/v10"
)

// NewMessageRoutes -.
func NewMessageRoutes(messageUseCase usecase.MessageUseCase, l logger.Interface) map[string]rabbitmq.Handler {
	r := &V1{
		messageUseCase: messageUseCase,
		l:              l,
		v:              validator.New(validator.WithRequiredStructEnabled()),
	}
	routes := make(map[string]rabbitmq.Handler)
	routes["user.deleted"] = r.userDeleted
	routes["room.deleted"] = r.roomDeleted
	routes["message.created"] = r.messageCreated
	routes["message.deleted"] = r.messageDeleted
	return routes
}
