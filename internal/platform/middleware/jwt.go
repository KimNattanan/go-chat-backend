package middleware

import (
	"errors"
	"net/http"

	authPb "github.com/KimNattanan/go-chat-backend/internal/auth/proto/v1"
	"github.com/KimNattanan/go-chat-backend/internal/platform/config"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/KimNattanan/go-chat-backend/pkg/requestid"
	"github.com/KimNattanan/go-chat-backend/pkg/responses"
	"github.com/KimNattanan/go-chat-backend/pkg/token"
	"github.com/labstack/echo/v5"
	"google.golang.org/grpc/status"
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

func jwtLog(c *echo.Context, l logger.Interface) logger.Interface {
	if rid := requestid.FromContext(c.Request().Context()); rid != "" {
		return l.With(requestid.MetadataKey, rid)
	}
	return l
}

func JWTMiddleware(l logger.Interface, cfg *config.Config, jwtMaker *token.JWTMaker, authGrpcClient authPb.AuthServiceClient) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			log := jwtLog(c, l)

			accessToken, err := readCookie(c, "access-token")
			if err != nil {
				log.Error(err, "JWTMiddleware")
				return responses.ErrorResponseCustom(c, http.StatusInternalServerError, "failed to read cookie")
			}
			accessClaims, err := jwtMaker.VerifyToken(accessToken, token.TokenTypeAccess)
			if err == nil {
				c.Set("userID", accessClaims.ID)
				return next(c)
			}
			refreshToken, err := readCookie(c, "refresh-token")
			if err != nil {
				log.Error(err, "JWTMiddleware")
				return responses.ErrorResponseCustom(c, http.StatusInternalServerError, "failed to read cookie")
			}
			refreshClaims, err := jwtMaker.VerifyToken(refreshToken, token.TokenTypeRefresh)
			if err != nil {
				log.Warn(err, "JWTMiddleware")
				return responses.ErrorResponseCustom(c, http.StatusUnauthorized, "unauthorized")
			}

			resp, err := authGrpcClient.RefreshToken(c.Request().Context(), &authPb.RefreshTokenRequest{
				UserId:    refreshClaims.ID,
				SessionId: refreshClaims.RegisteredClaims.ID,
			})
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					log.Error(err, "JWTMiddleware")
					return responses.ErrorResponseCustom(c, http.StatusInternalServerError, "internal server error")
				}
				code := responses.GrpcToHttpStatus(st.Code())
				if code >= http.StatusInternalServerError {
					log.Error(err, "JWTMiddleware")
				} else {
					log.Warn(err, "JWTMiddleware")
				}
				return responses.ErrorResponseCustom(c, code, st.Message())
			}

			c.SetCookie(&http.Cookie{
				Name:     "access-token",
				Value:    resp.Tokens.AccessToken,
				Expires:  resp.Tokens.AccessTokenExpiresAt.AsTime(),
				Path:     "/",
				HttpOnly: true,
				Secure:   cfg.App.ENV == "production",
				SameSite: http.SameSiteLaxMode,
			})
			c.SetCookie(&http.Cookie{
				Name:     "refresh-token",
				Value:    resp.Tokens.RefreshToken,
				Expires:  resp.Tokens.RefreshTokenExpiresAt.AsTime(),
				Path:     "/",
				HttpOnly: true,
				Secure:   cfg.App.ENV == "production",
				SameSite: http.SameSiteLaxMode,
			})
			c.Set("userID", resp.UserId)
			return next(c)
		}
	}
}
