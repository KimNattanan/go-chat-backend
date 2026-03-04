package responses

import (
	"github.com/labstack/echo/v5"
)

type Message struct {
	Message string `json:"message" example:"message"`
}

func MessageResponse(c *echo.Context, code int, msg string) error {
	return c.JSON(code, Message{Message: msg})
}
