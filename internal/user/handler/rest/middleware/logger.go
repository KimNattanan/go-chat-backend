package middleware

import (
	"strconv"
	"strings"

	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/labstack/echo/v5"
)

func buildRequestMessage(c *echo.Context) string {
	var result strings.Builder

	result.WriteString(c.RealIP())
	result.WriteString(" - ")
	result.WriteString(c.Request().Method)
	result.WriteString(" ")
	result.WriteString(c.Request().RequestURI)
	result.WriteString(" - ")
	result.WriteString(strconv.Itoa(c.Request().Response.StatusCode))

	return result.String()
}

func Logger(l logger.Interface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			err := next(c)
			l.Info("%s", buildRequestMessage(c))
			return err
		}
	}
}
