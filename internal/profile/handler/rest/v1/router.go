package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/profile/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

// NewProfileRoutes -.
func NewProfileRoutes(apiPublicGroup *echo.Group, apiPrivateGroup *echo.Group, profileUseCase usecase.ProfileUseCase, l logger.Interface) {
	r := &V1{
		profileUseCase: profileUseCase,
		l:              l,
		v:              validator.New(validator.WithRequiredStructEnabled()),
	}

	// Public Routes

	profilePublicGroup := apiPublicGroup.Group("/profiles")
	{
		profilePublicGroup.GET("/:userID", r.findProfile)
	}

	// Private Routes

	profilePrivateGroup := apiPublicGroup.Group("/profiles/me")
	{
		profilePrivateGroup.POST("", r.createProfile)
		profilePublicGroup.GET("", r.getProfile)
		profilePrivateGroup.PATCH("", r.patchProfile)
		profilePrivateGroup.DELETE("", r.deleteProfile)
	}

}
