package rest

import (
	v1 "github.com/KimNattanan/go-chat-backend/internal/message/handler/rest/v1"
	"github.com/KimNattanan/go-chat-backend/internal/message/usecase"
	"github.com/KimNattanan/go-chat-backend/internal/platform/config"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/labstack/echo/v5"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// NewRouter -.
// Swagger spec:
// @title       Go Clean Template API
// @description Using a translation service as an example
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(e *echo.Echo, cfg *config.Config, messageUseCase usecase.MessageUseCase, l logger.Interface, jwtMiddleware echo.MiddlewareFunc) {
	// Swagger
	if cfg.Swagger.Enabled {
		e.GET("/message/swagger/*", echoSwagger.WrapHandler)
	}

	// Routers
	apiPrivateGroup := e.Group("/v1", jwtMiddleware)
	v1.NewMessageRoutes(apiPrivateGroup, messageUseCase, l)
}
