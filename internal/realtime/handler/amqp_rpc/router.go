package amqp_rpc

import (
	v1 "github.com/KimNattanan/go-chat-backend/internal/realtime/handler/amqp_rpc/v1"
	"github.com/KimNattanan/go-chat-backend/internal/platform/wsserver"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/KimNattanan/go-chat-backend/pkg/rabbitmq"
)

// NewRouter -.
func NewRouter(wsServer *wsserver.Server, l logger.Interface) map[string]rabbitmq.Handler {
	return v1.NewRouter(wsServer, l)
}
