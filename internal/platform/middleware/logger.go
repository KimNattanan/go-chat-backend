package middleware

import (
	"strings"

	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/KimNattanan/go-chat-backend/pkg/requestid"
	"github.com/labstack/echo/v5"
)

func buildRequestMessage(c *echo.Context) string {
	var result strings.Builder

	result.WriteString(c.RealIP())
	result.WriteString(" - ")
	result.WriteString(c.Request().Method)
	result.WriteString(" ")
	result.WriteString(c.Request().RequestURI)

	return result.String()
}

func Logger(l logger.Interface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			err := next(c)
			log := l
			if rid, ok := c.Get(requestid.EchoContextKey).(string); ok && rid != "" {
				log = l.With(requestid.MetadataKey, rid)
			}
			log.Info("%s", buildRequestMessage(c))
			return err
		}
	}
}
