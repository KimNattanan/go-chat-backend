package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/profile/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

// NewProfileRoutes -.
func NewProfileRoutes(apiPublicGroup *echo.Group, profileUseCase usecase.ProfileUseCase, l logger.Interface) {
	r := &V1{
		profileUseCase: profileUseCase,
		l:              l,
		v:              validator.New(validator.WithRequiredStructEnabled()),
	}

	// Public Routes

	profilePublicGroup := apiPublicGroup.Group("/profiles")
	{
		profilePublicGroup.GET("/:id", r.findProfileByID)
		profilePublicGroup.PATCH("/:id", r.patchProfile)
	}

}
