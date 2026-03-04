package v1

import (
	"net/http"

	"github.com/KimNattanan/go-chat-backend/internal/profile/entity"
	"github.com/KimNattanan/go-chat-backend/internal/profile/handler/rest/v1/request"
	"github.com/KimNattanan/go-chat-backend/pkg/responses"
	"github.com/labstack/echo/v5"
)

func (r *V1) findProfileByID(c *echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Param("id")
	profile, err := r.profileUseCase.FindByID(ctx, userID)
	if err != nil {
		r.l.Error(err, "rest - v1 - findProfileByID")
		return responses.ErrorResponse(c, err)
	}
	return c.JSON(http.StatusOK, toProfileResponse(profile))
}

func (r *V1) patchProfile(c *echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Param("id")
	var req request.PatchProfileRequest
	if err := c.Bind(req); err != nil {
		r.l.Error(err, "rest - v1 - patchProfile")
		return responses.ErrorResponse(c, err)
	}
	profile := &entity.Profile{
		Name: req.Name,
	}
	profile, err := r.profileUseCase.Patch(ctx, userID, profile)
	if err != nil {
		r.l.Error(err, "rest - v1 - patchProfile")
		return responses.ErrorResponse(c, err)
	}
	return c.JSON(http.StatusOK, profile)
}
