package responses

import (
	"github.com/KimNattanan/go-chat-backend/pkg/apperror"
	"github.com/labstack/echo/v5"
)

type Error struct {
	Error string `json:"error" example:"message"`
}

func ErrorResponse(c *echo.Context, err error) error {
	code, errResp := apperror.Parse(err)
	return ErrorResponseCustom(c, code, errResp.Message)
}

func ErrorResponseCustom(c *echo.Context, code int, msg string) error {
	return c.JSON(code, Error{Error: msg})
}
