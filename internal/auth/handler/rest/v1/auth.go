package v1

import (
	"net/http"
	"time"

	"github.com/KimNattanan/go-chat-backend/internal/auth/handler/rest/v1/request"
	"github.com/KimNattanan/go-chat-backend/pkg/responses"
	"github.com/labstack/echo/v5"
)

func (r *V1) login(c *echo.Context) error {
	ctx := c.Request().Context()
	var req request.LoginRequest
	if err := c.Bind(&req); err != nil {
		r.l.Error(err, "rest - v1 - login")
		return responses.ErrorResponse(c, err)
	}
	_, accessToken, accessClaims, refreshToken, refreshClaims, err := r.authUseCase.Login(ctx, req.Email, req.Password)
	if err != nil {
		r.l.Error(err, "rest - v1 - login")
		return responses.ErrorResponse(c, err)
	}
	c.SetCookie(&http.Cookie{
		Name:     "access-token",
		Value:    accessToken,
		Expires:  accessClaims.ExpiresAt.Time,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.appEnv == "production",
		SameSite: http.SameSiteLaxMode,
	})
	c.SetCookie(&http.Cookie{
		Name:     "refresh-token",
		Value:    refreshToken,
		Expires:  refreshClaims.ExpiresAt.Time,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.appEnv == "production",
		SameSite: http.SameSiteLaxMode,
	})
	return responses.MessageResponse(c, http.StatusOK, "logged in successfully")
}

func (r *V1) register(c *echo.Context) error {
	ctx := c.Request().Context()
	var req request.RegisterRequest
	if err := c.Bind(&req); err != nil {
		r.l.Error(err, "rest - v1 - register")
		return responses.ErrorResponse(c, err)
	}
	_, accessToken, accessClaims, refreshToken, refreshClaims, err := r.authUseCase.Register(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		r.l.Error(err, "rest - v1 - register")
		return responses.ErrorResponse(c, err)
	}
	c.SetCookie(&http.Cookie{
		Name:     "access-token",
		Value:    accessToken,
		Expires:  accessClaims.ExpiresAt.Time,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.appEnv == "production",
		SameSite: http.SameSiteLaxMode,
	})
	c.SetCookie(&http.Cookie{
		Name:     "refresh-token",
		Value:    refreshToken,
		Expires:  refreshClaims.ExpiresAt.Time,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.appEnv == "production",
		SameSite: http.SameSiteLaxMode,
	})
	return responses.MessageResponse(c, http.StatusOK, "registered successfully")
}

func (r *V1) logout(c *echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     "access-token",
		Value:    "",
		Expires:  time.Now(),
		Path:     "/",
		HttpOnly: true,
		Secure:   r.appEnv == "production",
		SameSite: http.SameSiteLaxMode,
	})
	c.SetCookie(&http.Cookie{
		Name:     "refresh-token",
		Value:    "",
		Expires:  time.Now(),
		Path:     "/",
		HttpOnly: true,
		Secure:   r.appEnv == "production",
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}

func (r *V1) getUser(c *echo.Context) error {
	ctx := c.Request().Context()
	id := c.Get("userID").(string)
	user, err := r.authUseCase.FindUserByID(ctx, id)
	if err != nil {
		r.l.Error(err, "rest - v1 - getUser")
		return responses.ErrorResponse(c, err)
	}
	return c.JSON(http.StatusOK, toUserResponse(user))
}

func (r *V1) findUserByID(c *echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	user, err := r.authUseCase.FindUserByID(ctx, id)
	if err != nil {
		r.l.Error(err, "rest - v1 - findUserByID")
		return responses.ErrorResponse(c, err)
	}
	return c.JSON(http.StatusOK, toUserResponse(user))
}

func (r *V1) findUserByEmail(c *echo.Context) error {
	ctx := c.Request().Context()
	email := c.Param("email")
	user, err := r.authUseCase.FindUserByEmail(ctx, email)
	if err != nil {
		r.l.Error(err, "rest - v1 - findUserByEmail")
		return responses.ErrorResponse(c, err)
	}
	return c.JSON(http.StatusOK, toUserResponse(user))
}

func (r *V1) deleteUser(c *echo.Context) error {
	ctx := c.Request().Context()
	id := c.Get("userID").(string)
	if err := r.authUseCase.DeleteUser(ctx, id); err != nil {
		r.l.Error(err, "rest - v1 - deleteUser")
		return responses.ErrorResponse(c, err)
	}
	return responses.MessageResponse(c, http.StatusOK, "user deleted")
}
