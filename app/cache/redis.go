package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/config"
)

func NewRedisClient(cfg config.RedisConfig, log *logger.Logger) *redis.Client {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})
	_, err := rdb.Ping(ctx).Result()
	log.FatalIfErr(err, "[Redis] Failed To Connect Redis Client")
	log.Info().Msg("[Redis] Successfully Connected")
	return rdb
}
