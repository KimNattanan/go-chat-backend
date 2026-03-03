package rest

import (
	v1 "github.com/KimNattanan/go-chat-backend/internal/auth/handler/rest/v1"
	"github.com/KimNattanan/go-chat-backend/internal/auth/usecase"
	"github.com/KimNattanan/go-chat-backend/internal/platform/config"
	"github.com/KimNattanan/go-chat-backend/internal/platform/middleware"
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
func NewRouter(e *echo.Echo, cfg *config.Config, userUseCase usecase.UserUseCase, sessionUseCase usecase.SessionUseCase, l logger.Interface, jwtMiddleware echo.MiddlewareFunc) {
	// Options
	e.Use(middleware.Logger(l))
	e.Use(middleware.Recovery(l))

	// Swagger
	if cfg.Swagger.Enabled {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	// Routers
	apiV1Group := e.Group("/v1")
	{
		v1.NewAuthRoutes(apiV1Group, userUseCase, sessionUseCase, l, jwtMiddleware)
	}
}
