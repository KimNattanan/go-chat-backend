package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/labstack/echo/v5"
)

func buildPanicMessage(c *echo.Context, err any) string {
	var result strings.Builder

	result.WriteString(c.RealIP())
	result.WriteString(" - ")
	result.WriteString(c.Request().Method)
	result.WriteString(" ")
	result.WriteString(c.Request().RequestURI)
	result.WriteString(" PANIC DETECTED: ")
	result.WriteString(fmt.Sprintf("%v\n%s\n", err, debug.Stack())) //nolint: staticcheck,gocritic // it's okay for panic

	return result.String()
}

func Recovery(l logger.Interface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					l.Error(buildPanicMessage(c, r))
					err = echo.NewHTTPError(
						http.StatusInternalServerError,
						http.StatusText(http.StatusInternalServerError),
					)
				}
			}()
			return next(c)
		}
	}
}
