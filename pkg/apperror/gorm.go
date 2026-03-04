package apperror

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

func parseGormError(err error) (int, string, bool) {

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return http.StatusNotFound, "not found", true

	case errors.Is(err, gorm.ErrInvalidTransaction):
		return http.StatusBadRequest, "invalid transaction", true

	case errors.Is(err, gorm.ErrNotImplemented):
		return http.StatusNotImplemented, "feature not implemented", true

	case errors.Is(err, gorm.ErrMissingWhereClause):
		return http.StatusBadRequest, "missing where clause", true

	case errors.Is(err, gorm.ErrUnsupportedRelation):
		return http.StatusBadRequest, "unsupported relation", true

	case errors.Is(err, gorm.ErrPrimaryKeyRequired):
		return http.StatusBadRequest, "primary key required", true

	case errors.Is(err, gorm.ErrModelValueRequired):
		return http.StatusBadRequest, "model value required", true

	case errors.Is(err, gorm.ErrInvalidData):
		return http.StatusBadRequest, "invalid data", true

	case errors.Is(err, gorm.ErrUnsupportedDriver):
		return http.StatusInternalServerError, "unsupported driver", true

	case errors.Is(err, gorm.ErrRegistered):
		return http.StatusInternalServerError, "already registered", true

	case errors.Is(err, gorm.ErrInvalidField):
		return http.StatusBadRequest, "invalid field", true

	case errors.Is(err, gorm.ErrEmptySlice):
		return http.StatusBadRequest, "empty slice", true

	case errors.Is(err, gorm.ErrDryRunModeUnsupported):
		return http.StatusBadRequest, "dry run mode unsupported", true

	case errors.Is(err, gorm.ErrInvalidDB):
		return http.StatusInternalServerError, "invalid database instance", true

	case errors.Is(err, gorm.ErrInvalidValue):
		return http.StatusBadRequest, "invalid value", true

	case errors.Is(err, gorm.ErrDuplicatedKey):
		return http.StatusConflict, "duplicate key", true

	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return http.StatusConflict, "foreign key violated", true
	}

	return 0, "", false
}
