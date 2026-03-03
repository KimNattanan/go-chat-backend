package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/auth/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

// NewProfileRoutes -.
func NewAuthRoutes(apiV1Group *echo.Group, userUseCase usecase.UserUseCase, sessionUseCase usecase.SessionUseCase, l logger.Interface, jwtMiddleware echo.MiddlewareFunc) {
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
	}
	userGroup := apiV1Group.Group("/users")
	{
		userGroup.GET("/:id", r.findUserByID)
		userGroup.GET("/email/:email", r.findUserByEmail)
	}

	apiV1ProtectedGroup := apiV1Group.Group("")
	apiV1ProtectedGroup.Use(jwtMiddleware)

	authProtectedGroup := apiV1ProtectedGroup.Group("/auth")
	{
		authProtectedGroup.GET("/me", r.getUser)
		authProtectedGroup.POST("/logout", r.logout)
		authProtectedGroup.DELETE("/me", r.deleteUser)
	}
}
