package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/platform/wsserver"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/KimNattanan/go-chat-backend/pkg/rabbitmq"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

// NewWsRoutes -.
func NewWsRoutes(apiPrivateGroup *echo.Group, wsServer *wsserver.Server, mqPublisher rabbitmq.Publisher, l logger.Interface) {
	r := &V1{
		wsServer:    wsServer,
		mqPublisher: mqPublisher,
		l:           l,
		v:           validator.New(validator.WithRequiredStructEnabled()),
	}

	apiPrivateGroup.GET("/rooms/:roomID/ws", r.roomWebSocket)
}
