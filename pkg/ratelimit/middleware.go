package ratelimit

import (
	"net/http"

	"github.com/KimNattanan/go-chat-backend/pkg/responses"
	"github.com/labstack/echo/v5"
)

func RateLimitMiddleware(rl *RateLimiter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ip := c.RealIP()
			if !rl.Allow(ip) {
				return responses.ErrorResponseCustom(c, http.StatusTooManyRequests, "too many requests")
			}
			return next(c)
		}
	}
}
