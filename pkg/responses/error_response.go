package responses

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"google.golang.org/grpc/codes"
)

type Error struct {
	Error string `json:"error" example:"message"`
}

func ErrorResponse(c *echo.Context, code int, msg string) error {
	return c.JSON(code, Error{Error: msg})
}

func GrpcToHttpStatus(code codes.Code) int {
	switch code {

	case codes.OK:
		return http.StatusOK // 200

	case codes.InvalidArgument:
		return http.StatusBadRequest // 400

	case codes.FailedPrecondition:
		return http.StatusBadRequest // 400

	case codes.OutOfRange:
		return http.StatusBadRequest // 400

	case codes.Unauthenticated:
		return http.StatusUnauthorized // 401

	case codes.PermissionDenied:
		return http.StatusForbidden // 403

	case codes.NotFound:
		return http.StatusNotFound // 404

	case codes.AlreadyExists:
		return http.StatusConflict // 409

	case codes.Aborted:
		return http.StatusConflict // 409

	case codes.ResourceExhausted:
		return http.StatusTooManyRequests // 429

	case codes.Canceled:
		return 499 // Client closed request

	case codes.Internal:
		return http.StatusInternalServerError // 500

	case codes.Unknown:
		return http.StatusInternalServerError // 500

	case codes.DataLoss:
		return http.StatusInternalServerError // 500

	case codes.Unimplemented:
		return http.StatusNotImplemented // 501

	case codes.Unavailable:
		return http.StatusServiceUnavailable // 503

	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout // 504

	default:
		return http.StatusInternalServerError // fallback
	}
}
