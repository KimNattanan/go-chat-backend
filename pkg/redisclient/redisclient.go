package redisclient

import "github.com/redis/go-redis/v9"

func New(address, password string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	return rdb
}
