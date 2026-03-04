package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/auth/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

// NewAuthRoutes -.
func NewAuthRoutes(apiPublicGroup, apiPrivateGroup *echo.Group, authUseCase usecase.AuthUseCase, l logger.Interface) {
	r := &V1{
		authUseCase: authUseCase,
		l:           l,
		v:           validator.New(validator.WithRequiredStructEnabled()),
	}

	// Public Routes

	authPublicGroup := apiPublicGroup.Group("/auth")
	{
		authPublicGroup.POST("/login", r.login)
		authPublicGroup.POST("/register", r.register)
	}
	userPublicGroup := apiPublicGroup.Group("/users")
	{
		userPublicGroup.GET("/:id", r.findUserByID)
		userPublicGroup.GET("/email/:email", r.findUserByEmail)
	}

	// Private Routes

	authPrivateGroup := apiPrivateGroup.Group("/auth")
	{
		authPrivateGroup.GET("/me", r.getUser)
		authPrivateGroup.POST("/logout", r.logout)
		authPrivateGroup.DELETE("/me", r.deleteUser)
	}
}
