package service

import (
	"github.com/zhayt/cert-tz/storage"
	"go.uber.org/zap"
)

type Service struct {
	User IUserService
}

func NewService(storage *storage.Storage, l *zap.Logger) *Service {
	userService := NewUserService(storage, l)
	return &Service{
		User: userService,
	}
}
