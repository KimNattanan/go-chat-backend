package middleware

import (
	"errors"
	"net/http"

	authPb "github.com/KimNattanan/go-chat-backend/internal/auth/proto/v1"
	"github.com/KimNattanan/go-chat-backend/internal/platform/config"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
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

func JWTMiddleware(l logger.Interface, cfg *config.Config, jwtMaker *token.JWTMaker, authGrpcClient authPb.AuthServiceClient) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			accessToken, err := readCookie(c, "access-token")
			if err != nil {
				l.Error(err, "JWTMiddleware")
				return responses.ErrorResponseCustom(c, http.StatusInternalServerError, "failed to read cookie")
			}
			accessClaims, err := jwtMaker.VerfiyToken(accessToken)
			if err == nil {
				c.Set("userID", accessClaims.ID)
				return next(c)
			}
			refreshToken, err := readCookie(c, "refresh-token")
			if err != nil {
				l.Error(err, "JWTMiddleware")
				return responses.ErrorResponseCustom(c, http.StatusInternalServerError, "failed to read cookie")
			}
			refreshClaims, err := jwtMaker.VerfiyToken(refreshToken)
			if err != nil {
				l.Error(err, "JWTMiddleware")
				return responses.ErrorResponseCustom(c, http.StatusUnauthorized, "unauthorized")
			}

			// newAccessToken, newAccessClaims, err := jwtMaker.CreateToken(refreshClaims.ID, time.Second*time.Duration(cfg.JWT.AccessTTL))
			// if err != nil {
			// 	l.Error(err, "JWTMiddleware")
			// 	return responses.ErrorResponseCustom(c, http.StatusInternalServerError, "failed to create token")
			// }
			// newRefreshToken, newRefreshClaims, err := jwtMaker.CreateToken(refreshClaims.ID, time.Second*time.Duration(cfg.JWT.RefreshTTL))
			// if err != nil {
			// 	l.Error(err, "JWTMiddleware")
			// 	return responses.ErrorResponseCustom(c, http.StatusInternalServerError, "failed to create token")
			// }

			// grpc ->
			// 	 u := FindUserByID(UserID)
			//   s := FindSessionByID(SessionID)
			//   RevokeSession(SessionID)
			//   CreateSession(NewSessionID, u.ID, s.CreatedAt, ExpiresAt)
			resp, err := authGrpcClient.RefreshToken(c.Request().Context(), &authPb.RefreshTokenRequest{
				UserId:    refreshClaims.ID,
				SessionId: refreshClaims.RegisteredClaims.ID,
				// NewSessionId: newRefreshClaims.RegisteredClaims.ID,
				// ExpiresAt:    timestamppb.New(newRefreshClaims.RegisteredClaims.ExpiresAt.Time),
			})
			if err != nil {
				l.Error(err, "JWTMiddleware")
				st, ok := status.FromError(err)
				if !ok {
					return responses.ErrorResponseCustom(c, http.StatusInternalServerError, "internal server error")
				}
				return responses.ErrorResponseCustom(c, responses.GrpcToHttpStatus(st.Code()), st.Message())
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
