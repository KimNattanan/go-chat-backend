package v1

import (
	"github.com/KimNattanan/go-chat-backend/pkg/responses"
	"github.com/labstack/echo/v5"
)

func errorResponse(c *echo.Context, code int, msg string) error {
	return c.JSON(code, responses.ErrorResponse{Error: msg})
}
