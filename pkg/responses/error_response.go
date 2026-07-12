package responses

import (
	"net/http"

	"github.com/KimNattanan/go-chat-backend/pkg/apperror"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/KimNattanan/go-chat-backend/pkg/requestid"
	"github.com/labstack/echo/v5"
)

type Error struct {
	Error  string            `json:"error" example:"message"`
	Code   string            `json:"code,omitempty" example:"UNAUTHORIZED"`
	Errors map[string]string `json:"errors,omitempty"`
}

func ErrorResponse(c *echo.Context, err error) error {
	code, errResp := apperror.ParseHttp(err)
	return c.JSON(code, Error{
		Error:  errResp.Message,
		Code:   errResp.Code,
		Errors: errResp.Errors,
	})
}

func ErrorResponseCustom(c *echo.Context, code int, msg string) error {
	return c.JSON(code, Error{
		Error: msg,
		Code:  apperror.CodeFromHTTPStatus(code),
	})
}

// LogAndErrorResponse logs at Warn for expected client errors (<500) and Error for server errors,
// then returns the standard JSON error body.
func LogAndErrorResponse(c *echo.Context, l logger.Interface, err error, op string) error {
	code, _ := apperror.ParseHttp(err)
	log := l
	if rid := requestid.FromContext(c.Request().Context()); rid != "" {
		log = l.With(requestid.MetadataKey, rid)
	}
	if code >= http.StatusInternalServerError {
		log.Error(err, op)
	} else {
		log.Warn(err, op)
	}
	return ErrorResponse(c, err)
}

// LogAMQP logs AMQP handler errors: Warn for client/validation-style failures, Error otherwise.
func LogAMQP(l logger.Interface, err error, op string) {
	code, _ := apperror.ParseHttp(err)
	if code >= http.StatusInternalServerError {
		l.Error(err, op)
	} else {
		l.Warn(err, op)
	}
}
