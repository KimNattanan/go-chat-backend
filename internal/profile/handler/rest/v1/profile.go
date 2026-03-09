package v1

import (
	"net/http"

	"github.com/KimNattanan/go-chat-backend/internal/profile/entity"
	"github.com/KimNattanan/go-chat-backend/internal/profile/handler/rest/v1/request"
	"github.com/KimNattanan/go-chat-backend/pkg/responses"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func (r *V1) findProfile(c *echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Param("userID")

	profile, err := r.profileUseCase.FindByUserID(ctx, userID)
	if err != nil {
		r.l.Error(err, "rest - v1 - findProfile")
		return responses.ErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, toProfileResponse(profile))
}

func (r *V1) getProfile(c *echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Get("userID").(string)

	profile, err := r.profileUseCase.FindByUserID(ctx, userID)
	if err != nil {
		r.l.Error(err, "rest - v1 - getProfile")
		return responses.ErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, toProfileResponse(profile))
}

func (r *V1) createProfile(c *echo.Context) error {
	ctx := c.Request().Context()

	userIDStr := c.Get("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		r.l.Error(err, "rest - v1 - createProfile")
		return responses.ErrorResponse(c, err)
	}

	var req request.CreateProfileRequest
	if err := c.Bind(&req); err != nil {
		r.l.Error(err, "rest - v1 - createProfile")
		return responses.ErrorResponse(c, err)
	}
	if err := r.v.Struct(&req); err != nil {
		r.l.Error(err, "rest - v1 - createProfile")
		return responses.ErrorResponse(c, err)
	}

	profile := &entity.Profile{
		UserID: userID,
		Email:  req.Email,
		Name:   req.Name,
	}
	if err := r.profileUseCase.Create(ctx, profile); err != nil {
		r.l.Error(err, "rest - v1 - createProfile")
		return responses.ErrorResponse(c, err)
	}

	return c.JSON(http.StatusCreated, toProfileResponse(profile))
}

func (r *V1) patchProfile(c *echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Get("userID").(string)

	var req request.PatchProfileRequest
	if err := c.Bind(&req); err != nil {
		r.l.Error(err, "rest - v1 - patchProfile")
		return responses.ErrorResponse(c, err)
	}
	if err := r.v.Struct(&req); err != nil {
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

func (r *V1) deleteProfile(c *echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Get("userID").(string)

	if err := r.profileUseCase.Delete(ctx, userID); err != nil {
		r.l.Error(err, "rest - v1 - deleteProfile")
		return responses.ErrorResponse(c, err)
	}

	return responses.MessageResponse(c, http.StatusOK, "profile deleted")
}
