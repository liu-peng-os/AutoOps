package redis

import (
	"context"

	"dodevops-api/common/config"
	"github.com/go-redis/redis/v8"
)

var (
	RedisDb *redis.Client
)

func SetupRedisDb() error {
	ctx := context.Background()
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Address,
		Password: config.Config.Redis.Password,
		DB:       0,
	})

	if _, err := RedisDb.Ping(ctx).Result(); err != nil {
		RedisDb = nil
		return err
	}

	return nil
}
