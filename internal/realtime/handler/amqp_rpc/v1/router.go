package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/platform/wsserver"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/KimNattanan/go-chat-backend/pkg/rabbitmq"
	"github.com/go-playground/validator/v10"
)

// NewRouter -.
func NewRouter(wsServer *wsserver.Server, l logger.Interface) map[string]rabbitmq.Handler {
	r := &V1{
		wsServer: wsServer,
		l:        l,
		v:        validator.New(validator.WithRequiredStructEnabled()),
	}
	routes := make(map[string]rabbitmq.Handler)
	routes["message.created"] = r.messageCreated
	routes["message.deleted"] = r.messageDeleted
	return routes
}
