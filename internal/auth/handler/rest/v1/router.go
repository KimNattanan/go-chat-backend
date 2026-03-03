package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/auth/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

// NewProfileRoutes -.
func NewAuthRoutes(apiV1Group *echo.Group, userUseCase usecase.UserUseCase, sessionUseCase usecase.SessionUseCase, l logger.Interface) {
	r := &V1{
		userUseCase:    userUseCase,
		sessionUseCase: sessionUseCase,
		l:              l,
		v:              validator.New(validator.WithRequiredStructEnabled()),
	}

	authGroup := apiV1Group.Group("/auth")
	{
		authGroup.POST("/login", r.login)
		authGroup.POST("/register", r.register)
		authGroup.POST("/logout", r.logout)
	}
	userGroup := apiV1Group.Group("/users")
	{
		userGroup.GET("/:id", r.findUserByID)
		userGroup.GET("/email/:email", r.findUserByEmail)
		userGroup.DELETE("/:id", r.deleteUser)
	}

}
