package http

import (
	"errors"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func ParseRedisError(err error) (int, string, bool) {
	if err == nil {
		return 0, "", false
	}

	if errors.Is(err, redis.Nil) {
		return http.StatusNotFound, "not found", true
	}

	var redisErr redis.Error
	if errors.As(err, &redisErr) {
		return http.StatusInternalServerError, "redis error", true
	}

	return 0, "", false
}
