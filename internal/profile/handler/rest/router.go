package rest

import (
	"github.com/KimNattanan/go-chat-backend/internal/platform/config"
	"github.com/KimNattanan/go-chat-backend/internal/platform/middleware"
	v1 "github.com/KimNattanan/go-chat-backend/internal/profile/handler/rest/v1"
	"github.com/KimNattanan/go-chat-backend/internal/profile/usecase"
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
func NewRouter(e *echo.Echo, cfg *config.Config, profileUseCase usecase.ProfileUseCase, l logger.Interface) {
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
		v1.NewProfileRoutes(apiV1Group, profileUseCase, l)
	}
}
