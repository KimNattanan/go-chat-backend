package ws

import (
	"github.com/KimNattanan/go-chat-backend/internal/platform/wsserver"
	v1 "github.com/KimNattanan/go-chat-backend/internal/realtime/handler/ws/v1"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/KimNattanan/go-chat-backend/pkg/rabbitmq"
	"github.com/labstack/echo/v5"
)

func NewRouter(e *echo.Echo, wsServer *wsserver.Server, amqpPublisher *rabbitmq.Publisher, l logger.Interface, jwtMiddleware echo.MiddlewareFunc) {
	apiPrivateGroup := e.Group("/v1", jwtMiddleware)
	v1.NewWsRoutes(apiPrivateGroup, wsServer, amqpPublisher, l)
}
