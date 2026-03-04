package apperror

import (
	"errors"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func parseRedisError(err error) (int, string, bool) {
	if err == nil {
		return 0, "", false
	}

	// Nil (key not found)
	if errors.Is(err, redis.Nil) {
		return http.StatusNotFound, "not found", true
	}

	// Redis specific error types
	var redisErr redis.Error
	if errors.As(err, &redisErr) {
		// Most Redis server errors (WRONGTYPE, NOAUTH, etc.)
		return http.StatusBadRequest, redisErr.Error(), true
	}

	return 0, "", false
}
