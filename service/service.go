package service

import (
	"github.com/zhayt/cert-tz/storage"
	"go.uber.org/zap"
)

type Service struct {
	User     IUserService
	SubStr   ISubStrService
	Analysis IAnalysisService
	Counter  ICounterService
}

func NewService(storage *storage.Storage, l *zap.Logger) *Service {
	validateService := NewValidateService()

	userService := NewUserService(storage, validateService, l)
	subStrService := NewSubStrService(validateService, l)
	analysisService := NewAnalysisService(l)
	counterService := NewCounterService(storage, l)

	return &Service{
		User:     userService,
		SubStr:   subStrService,
		Analysis: analysisService,
		Counter:  counterService,
	}
}
