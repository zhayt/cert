package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type CounterStorage struct {
	client *redis.Client
}

func (r *CounterStorage) IncreaseCounter(ctx context.Context, key string, val int64) error {
	err := r.client.IncrBy(ctx, key, val).Err()
	if err != nil {
		return fmt.Errorf("cannot increase counter: %w", err)
	}

	return nil
}

func (r *CounterStorage) DecreaseCounter(ctx context.Context, key string, val int64) error {
	err := r.client.DecrBy(ctx, key, val).Err()
	if err != nil {
		return fmt.Errorf("cannot decrease counter: %w", err)
	}

	return nil
}

func (r *CounterStorage) GetCounter(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("cannot get counter: %w", err)
	}

	return val, nil
}

func NewCounterStorage(client *redis.Client) *CounterStorage {
	return &CounterStorage{client: client}
}
