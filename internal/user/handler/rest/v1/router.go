package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/user/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

// NewUserRoutes -.
func NewUserRoutes(apiV1Group *echo.Group, profileUseCase usecase.ProfileUseCase, l logger.Interface) {
	r := &V1{
		profileUseCase: profileUseCase,
		l:              l,
		v:              validator.New(validator.WithRequiredStructEnabled()),
	}

	profileGroup := apiV1Group.Group("/profiles")
	{
		profileGroup.GET("/:id", r.findProfileByID)
		profileGroup.PATCH("/:id", r.patchProfile)
	}
}
