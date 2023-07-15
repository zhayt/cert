package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/zhayt/cert-tz/config"
	"log"
	"net"
	"time"
)

func NewRedisClient(cfg *config.Config) (*redis.Client, error) {
	redisHost := net.JoinHostPort(cfg.RedisHost, cfg.RedisPort)

	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: cfg.RedisPassword, // no password set
		DB:       cfg.RedisDB,       // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("cannot connect to redis: %w", err)
	}

	log.Println(pong)

	return client, nil
}
