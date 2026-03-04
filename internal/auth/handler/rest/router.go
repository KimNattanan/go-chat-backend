package rest

import (
	v1 "github.com/KimNattanan/go-chat-backend/internal/auth/handler/rest/v1"
	"github.com/KimNattanan/go-chat-backend/internal/auth/usecase"
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
func NewRouter(e *echo.Echo, cfg *config.Config, authUseCase usecase.AuthUseCase, l logger.Interface, jwtMiddleware echo.MiddlewareFunc) {
	// Swagger
	if cfg.Swagger.Enabled {
		e.GET("/auth/swagger/*", echoSwagger.WrapHandler)
	}

	// Routers
	apiPublicGroup := e.Group("/v1")
	apiPrivateGroup := e.Group("/v1", jwtMiddleware)
	v1.NewAuthRoutes(apiPublicGroup, apiPrivateGroup, authUseCase, l, cfg.App.ENV)
}
