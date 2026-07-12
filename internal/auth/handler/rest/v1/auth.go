package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/KimNattanan/go-chat-backend/internal/auth/handler/rest/v1/request"
	"github.com/KimNattanan/go-chat-backend/pkg/responses"
	"github.com/labstack/echo/v5"
)

func readCookie(c *echo.Context, name string) (string, error) {
	cookie, err := c.Cookie(name)
	if errors.Is(err, echo.ErrCookieNotFound) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (r *V1) login(c *echo.Context) error {
	ctx := c.Request().Context()

	var req request.LoginRequest
	if err := c.Bind(&req); err != nil {
		return responses.LogAndErrorResponse(c, r.l, err, "rest - v1 - login")
	}
	if err := r.v.Struct(&req); err != nil {
		return responses.LogAndErrorResponse(c, r.l, err, "rest - v1 - login")
	}

	result, err := r.authUseCase.Login(ctx, req.Email, req.Password)
	if err != nil {
		return responses.LogAndErrorResponse(c, r.l, err, "rest - v1 - login")
	}

	c.SetCookie(&http.Cookie{
		Name:     "access-token",
		Value:    result.AccessToken,
		Expires:  result.AccessClaims.ExpiresAt.Time,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.appEnv == "production",
		SameSite: http.SameSiteLaxMode,
	})
	c.SetCookie(&http.Cookie{
		Name:     "refresh-token",
		Value:    result.RefreshToken,
		Expires:  result.RefreshClaims.ExpiresAt.Time,
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
		return responses.LogAndErrorResponse(c, r.l, err, "rest - v1 - register")
	}
	if err := r.v.Struct(&req); err != nil {
		return responses.LogAndErrorResponse(c, r.l, err, "rest - v1 - register")
	}

	result, err := r.authUseCase.Register(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		return responses.LogAndErrorResponse(c, r.l, err, "rest - v1 - register")
	}

	c.SetCookie(&http.Cookie{
		Name:     "access-token",
		Value:    result.AccessToken,
		Expires:  result.AccessClaims.ExpiresAt.Time,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.appEnv == "production",
		SameSite: http.SameSiteLaxMode,
	})
	c.SetCookie(&http.Cookie{
		Name:     "refresh-token",
		Value:    result.RefreshToken,
		Expires:  result.RefreshClaims.ExpiresAt.Time,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.appEnv == "production",
		SameSite: http.SameSiteLaxMode,
	})

	return responses.MessageResponse(c, http.StatusOK, "registered successfully")
}

func clearAuthCookies(c *echo.Context, appEnv string) {
	c.SetCookie(&http.Cookie{
		Name:     "access-token",
		Value:    "",
		Expires:  time.Now(),
		Path:     "/",
		HttpOnly: true,
		Secure:   appEnv == "production",
		SameSite: http.SameSiteLaxMode,
	})
	c.SetCookie(&http.Cookie{
		Name:     "refresh-token",
		Value:    "",
		Expires:  time.Now(),
		Path:     "/",
		HttpOnly: true,
		Secure:   appEnv == "production",
		SameSite: http.SameSiteLaxMode,
	})
}

func (r *V1) logout(c *echo.Context) error {
	ctx := c.Request().Context()

	refreshToken, err := readCookie(c, "refresh-token")
	if err != nil {
		return responses.LogAndErrorResponse(c, r.l, err, "rest - v1 - logout")
	}

	if err := r.authUseCase.Logout(ctx, refreshToken); err != nil {
		return responses.LogAndErrorResponse(c, r.l, err, "rest - v1 - logout")
	}

	clearAuthCookies(c, r.appEnv)
	return responses.MessageResponse(c, http.StatusOK, "logged out successfully")
}

func (r *V1) getUser(c *echo.Context) error {
	ctx := c.Request().Context()
	id := c.Get("userID").(string)

	user, err := r.authUseCase.FindUserByID(ctx, id)
	if err != nil {
		return responses.LogAndErrorResponse(c, r.l, err, "rest - v1 - getUser")
	}

	return c.JSON(http.StatusOK, toUserResponse(user))
}

func (r *V1) findUserByID(c *echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	user, err := r.authUseCase.FindUserByID(ctx, id)
	if err != nil {
		return responses.LogAndErrorResponse(c, r.l, err, "rest - v1 - findUserByID")
	}

	return c.JSON(http.StatusOK, toUserResponse(user))
}

func (r *V1) findUserByEmail(c *echo.Context) error {
	ctx := c.Request().Context()
	email := c.Param("email")

	user, err := r.authUseCase.FindUserByEmail(ctx, email)
	if err != nil {
		return responses.LogAndErrorResponse(c, r.l, err, "rest - v1 - findUserByEmail")
	}

	return c.JSON(http.StatusOK, toUserResponse(user))
}

func (r *V1) deleteUser(c *echo.Context) error {
	ctx := c.Request().Context()
	id := c.Get("userID").(string)

	if err := r.authUseCase.DeleteUser(ctx, id); err != nil {
		return responses.LogAndErrorResponse(c, r.l, err, "rest - v1 - deleteUser")
	}

	clearAuthCookies(c, r.appEnv)
	return responses.MessageResponse(c, http.StatusOK, "user deleted")
}
