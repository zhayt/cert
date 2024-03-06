package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/zhayt/cert-tz/internal/model"
	"github.com/zhayt/cert-tz/internal/storage"
	"github.com/zhayt/cert-tz/internal/storage/mocks"
	"go.uber.org/zap"
	"testing"
)

func TestUserService_CreateUser(t *testing.T) {
	type args struct {
		ctx  context.Context
		user model.User
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"success", args{context.Background(), model.User{FirstName: "Cert", LastName: "Test"}}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userStorage := mocks.NewIUserStorage(t)

			stor := &storage.Storage{UserStorage: userStorage}
			s := &UserService{
				storage:  stor,
				validate: NewValidateService(),
				l:        zap.NewNop(),
			}

			if tt.args.user.FirstName == "Cert" {
				userStorage.On("CreateUser", mock.Anything, mock.AnythingOfType("model.User")).
					Return(func(ctx context.Context, user model.User) (uint64, error) {
						return 1, nil
					}).Once()
			} else {
				userStorage.On("CreateUser", mock.Anything, mock.AnythingOfType("model.User")).
					Return(func(ctx context.Context, user model.User) (uint64, error) {
						return 0, errors.New("some error")
					}).Once()
			}

			got, err := s.CreateUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	type args struct {
		ctx  context.Context
		user model.User
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"success", args{context.Background(), model.User{FirstName: "Cert", LastName: "Updated"}}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userStorage := mocks.NewIUserStorage(t)

			stor := &storage.Storage{UserStorage: userStorage}
			s := &UserService{
				storage:  stor,
				validate: NewValidateService(),
				l:        zap.NewNop(),
			}
			if tt.args.user.FirstName == "Cert" {
				userStorage.On("UpdateUser", mock.Anything, mock.AnythingOfType("model.User")).
					Return(func(ctx context.Context, user model.User) (uint64, error) {
						return 1, nil
					}).Once()
			} else {
				userStorage.On("UpdateUser", mock.Anything, mock.AnythingOfType("model.User")).
					Return(func(ctx context.Context, user model.User) (uint64, error) {
						return 0, errors.New("some error")
					}).Once()
			}
			got, err := s.UpdateUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UpdateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
