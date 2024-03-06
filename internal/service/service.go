package service

import (
	"github.com/zhayt/cert-tz/internal/storage"
	"go.uber.org/zap"
)

type Service struct {
	User     IUserService
	SubStr   ISubStrService
	Analysis IAnalysisService
	Counter  ICounterService
	Hash     IHashService
}

func NewService(storage *storage.Storage, l *zap.Logger) *Service {
	validateService := NewValidateService()

	userService := NewUserService(storage, validateService, l)
	subStrService := NewSubStrService(validateService, l)
	analysisService := NewAnalysisService(l)
	counterService := NewCounterService(storage, l)
	hashService := NewHashService(storage, validateService, l)

	return &Service{
		User:     userService,
		SubStr:   subStrService,
		Analysis: analysisService,
		Counter:  counterService,
		Hash:     hashService,
	}
}
