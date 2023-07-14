package service

import (
	"github.com/zhayt/cert-tz/storage"
	"go.uber.org/zap"
)

type Service struct {
	User   IUserService
	SubStr ISubStrService
}

func NewService(storage *storage.Storage, l *zap.Logger) *Service {
	validateService := NewValidateService()
	userService := NewUserService(storage, validateService, l)
	subStrService := NewSubStrService(validateService, l)

	return &Service{
		User:   userService,
		SubStr: subStrService,
	}
}
