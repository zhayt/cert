package service

import (
	"context"
	"github.com/zhayt/cert-tz/storage"
	"go.uber.org/zap"
)

//go:generate mockery --name ICounterService
type ICounterService interface {
	IncreaseCounter(ctx context.Context, key string, val int64) error
	DecreaseCounter(ctx context.Context, key string, val int64) error
	GetCounter(ctx context.Context, key string) (string, error)
}

type CounterService struct {
	storage *storage.Storage
	l       *zap.Logger
}

func (s *CounterService) IncreaseCounter(ctx context.Context, key string, val int64) error {
	return s.storage.Cache.IncreaseCounter(ctx, key, val)
}

func (s *CounterService) DecreaseCounter(ctx context.Context, key string, val int64) error {
	return s.storage.Cache.DecreaseCounter(ctx, key, val)
}

func (s *CounterService) GetCounter(ctx context.Context, key string) (string, error) {
	return s.storage.Cache.GetCounter(ctx, key)
}

func NewCounterService(storage *storage.Storage, l *zap.Logger) *CounterService {
	return &CounterService{storage: storage, l: l}
}
