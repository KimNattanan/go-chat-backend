package apperror

import (
	"context"
	"errors"

	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/KimNattanan/go-chat-backend/pkg/requestid"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// ParseGrpcLogged logs the original error for ops, then maps it to a gRPC status.
func ParseGrpcLogged(ctx context.Context, l logger.Interface, err error, op string) error {
	if err == nil {
		return nil
	}
	log := l
	if rid := requestid.FromContext(ctx); rid != "" {
		log = l.With(requestid.MetadataKey, rid)
	}
	code, _ := ParseHttp(err)
	if code >= 500 {
		log.Error(err, op)
	} else {
		log.Warn(err, op)
	}
	return ParseGrpc(err)
}

func ParseGrpc(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := status.FromError(err); ok {
		return err // already a gRPC status
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return status.Error(httpToGrpcCode(appErr.Code), appErr.Message)
	}

	var redisErr redis.Error
	var validationErrs validator.ValidationErrors

	if errors.As(err, &validationErrs) {
		return status.Error(codes.InvalidArgument, "validation failed")
	}

	switch {

	// GORM

	case errors.Is(err, gorm.ErrRecordNotFound):
		return status.Error(codes.NotFound, "not found")

	case errors.Is(err, gorm.ErrInvalidTransaction):
		return status.Error(codes.InvalidArgument, "invalid transaction")

	case errors.Is(err, gorm.ErrNotImplemented):
		return status.Error(codes.Unimplemented, "feature not implemented")

	case errors.Is(err, gorm.ErrMissingWhereClause):
		return status.Error(codes.InvalidArgument, "missing where clause")

	case errors.Is(err, gorm.ErrUnsupportedRelation):
		return status.Error(codes.InvalidArgument, "unsupported relation")

	case errors.Is(err, gorm.ErrPrimaryKeyRequired):
		return status.Error(codes.InvalidArgument, "primary key required")

	case errors.Is(err, gorm.ErrModelValueRequired):
		return status.Error(codes.InvalidArgument, "model value required")

	case errors.Is(err, gorm.ErrInvalidData):
		return status.Error(codes.InvalidArgument, "invalid data")

	case errors.Is(err, gorm.ErrUnsupportedDriver):
		return status.Error(codes.Internal, "unsupported driver")

	case errors.Is(err, gorm.ErrRegistered):
		return status.Error(codes.Internal, "already registered")

	case errors.Is(err, gorm.ErrInvalidField):
		return status.Error(codes.InvalidArgument, "invalid field")

	case errors.Is(err, gorm.ErrEmptySlice):
		return status.Error(codes.InvalidArgument, "empty slice")

	case errors.Is(err, gorm.ErrDryRunModeUnsupported):
		return status.Error(codes.InvalidArgument, "dry run mode unsupported")

	case errors.Is(err, gorm.ErrInvalidDB):
		return status.Error(codes.Internal, "invalid database instance")

	case errors.Is(err, gorm.ErrInvalidValue):
		return status.Error(codes.InvalidArgument, "invalid value")

	case errors.Is(err, gorm.ErrDuplicatedKey):
		return status.Error(codes.AlreadyExists, "duplicate key")

	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return status.Error(codes.AlreadyExists, "foreign key violated")

	// Redis

	case errors.Is(err, redis.Nil):
		return status.Error(codes.NotFound, "not found")

	case errors.As(err, &redisErr):
		return status.Error(codes.Internal, "redis error")
	}

	return status.Error(codes.Internal, "internal server error")
}

func httpToGrpcCode(httpCode int) codes.Code {
	switch httpCode {
	case 400:
		return codes.InvalidArgument
	case 401:
		return codes.Unauthenticated
	case 403:
		return codes.PermissionDenied
	case 404:
		return codes.NotFound
	case 409:
		return codes.AlreadyExists
	default:
		return codes.Internal
	}
}
