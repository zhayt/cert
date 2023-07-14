package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/zhayt/cert-tz/model"
	"github.com/zhayt/cert-tz/storage/postgre"
)

type IUserStorage interface {
	CreateUser(ctx context.Context, user model.User) (uint64, error)
	GetUser(ctx context.Context, userID uint64) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) (uint64, error)
	DeleteUser(ctx context.Context, userID uint64) error
}

type Storage struct {
	UserStorage IUserStorage
}

func NewStorage(db *sqlx.DB) *Storage {
	userStorage := postgre.NewUserStorage(db)

	return &Storage{UserStorage: userStorage}
}
