package apperror

import (
	"errors"
	"net/http"

	httpAppError "github.com/KimNattanan/go-chat-backend/pkg/apperror/http"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Message string            `json:"message"`
	Code    string            `json:"code,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func ParseHttp(err error) (int, ErrorResponse) {
	if err == nil {
		return http.StatusOK, ErrorResponse{}
	}

	// Custom AppError (highest priority)
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code, ErrorResponse{
			Message: appErr.Message,
			Code:    CodeFromHTTPStatus(appErr.Code),
		}
	}

	// Validator Errors
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		return http.StatusBadRequest, ErrorResponse{
			Message: "validation failed",
			Code:    CodeValidationFailed,
			Errors:  httpAppError.ParseValidationErrors(validationErrs),
		}
	}

	// GORM Errors
	if code, msg, ok := httpAppError.ParseGormError(err); ok {
		return code, ErrorResponse{
			Message: msg,
			Code:    CodeFromHTTPStatus(code),
		}
	}

	// Redis Errors
	if code, msg, ok := httpAppError.ParseRedisError(err); ok {
		return code, ErrorResponse{
			Message: msg,
			Code:    CodeFromHTTPStatus(code),
		}
	}

	// Default fallback
	return http.StatusInternalServerError, ErrorResponse{
		Message: "internal server error",
		Code:    CodeInternal,
	}
}
