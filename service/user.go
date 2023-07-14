package service

import (
	"context"
	"errors"
	"github.com/zhayt/cert-tz/model"
	"github.com/zhayt/cert-tz/storage"
	"go.uber.org/zap"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

type IUserService interface {
	CreateUser(ctx context.Context, user model.User) (uint64, error)
	GetUser(ctx context.Context, userID uint64) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) (uint64, error)
	DeleteUser(ctx context.Context, userID uint64) error
}

var (
	ErrInvalidData = errors.New("invalid data")
)

type UserService struct {
	storage  *storage.Storage
	validate *ValidateService
	l        *zap.Logger
}

func (s *UserService) CreateUser(ctx context.Context, user model.User) (uint64, error) {
	if err := s.validate.validateStruct(user); err != nil {
		s.l.Error("validate error", zap.Error(err))
		return 0, ErrInvalidData
	}

	user.FirstName = cases.Title(language.Und, cases.NoLower).String(strings.ToLower(user.FirstName))
	user.LastName = cases.Title(language.Und, cases.NoLower).String(strings.ToLower(user.LastName))

	return s.storage.UserStorage.CreateUser(ctx, user)
}

func (s *UserService) GetUser(ctx context.Context, userID uint64) (model.User, error) {
	return s.storage.UserStorage.GetUser(ctx, userID)
}

func (s *UserService) UpdateUser(ctx context.Context, user model.User) (uint64, error) {
	if err := s.validate.validateStruct(user); err != nil {
		s.l.Error("validate error", zap.Error(err))
		return 0, ErrInvalidData
	}

	user.FirstName = cases.Title(language.Und, cases.NoLower).String(strings.ToLower(user.FirstName))
	user.LastName = cases.Title(language.Und, cases.NoLower).String(strings.ToLower(user.LastName))

	return s.storage.UserStorage.UpdateUser(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, userID uint64) error {
	return s.storage.UserStorage.DeleteUser(ctx, uint64(userID))
}

func NewUserService(storage *storage.Storage, validate *ValidateService, l *zap.Logger) *UserService {
	return &UserService{
		storage:  storage,
		validate: validate,
		l:        l,
	}
}
