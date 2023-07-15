package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/zhayt/cert-tz/model"
	"github.com/zhayt/cert-tz/storage/postgre"
	cache "github.com/zhayt/cert-tz/storage/redis"
)

type IUserStorage interface {
	CreateUser(ctx context.Context, user model.User) (uint64, error)
	GetUser(ctx context.Context, userID uint64) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) (uint64, error)
	DeleteUser(ctx context.Context, userID uint64) error
}

type ICounterStorage interface {
	IncreaseCounter(ctx context.Context, key string, val int64) error
	DecreaseCounter(ctx context.Context, key string, val int64) error
	GetCounter(ctx context.Context, key string) (string, error)
}

type Storage struct {
	UserStorage IUserStorage
	Cache       ICounterStorage
}

func NewStorage(db *sqlx.DB, redisClient *redis.Client) *Storage {
	userStorage := postgre.NewUserStorage(db)
	counterStorage := cache.NewCounterStorage(redisClient)

	return &Storage{
		UserStorage: userStorage,
		Cache:       counterStorage,
	}
}
