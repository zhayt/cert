package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/zhayt/cert-tz/internal/model"
	postgre2 "github.com/zhayt/cert-tz/internal/storage/postgre"
	cache "github.com/zhayt/cert-tz/internal/storage/redis"
)

//go:generate mockery --name IUserStorage
type IUserStorage interface {
	CreateUser(ctx context.Context, user model.User) (uint64, error)
	GetUser(ctx context.Context, userID uint64) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) (uint64, error)
	DeleteUser(ctx context.Context, userID uint64) error
}

//go:generate mockery --name ICounterStorage
type ICounterStorage interface {
	IncreaseCounter(ctx context.Context, key string, val int64) error
	DecreaseCounter(ctx context.Context, key string, val int64) error
	GetCounter(ctx context.Context, key string) (string, error)
}

//go:generate mockery --name IHashStorage
type IHashStorage interface {
	CreateHash(hash model.CertHash) (uint64, error)
	GetHash(hashID uint64) (model.CertHash, error)
	UpdateHash(hash model.CertHash) error
}

type Storage struct {
	UserStorage IUserStorage
	Cache       ICounterStorage
	Hash        IHashStorage
}

func NewStorage(db *sqlx.DB, redisClient *redis.Client) *Storage {
	userStorage := postgre2.NewUserStorage(db)
	counterStorage := cache.NewCounterStorage(redisClient)
	hashStorage := postgre2.NewHashStorage(db)

	return &Storage{
		UserStorage: userStorage,
		Cache:       counterStorage,
		Hash:        hashStorage,
	}
}
