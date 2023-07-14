package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/zhayt/cert-tz/model"
	"github.com/zhayt/cert-tz/storage"
	"go.uber.org/zap"
	"testing"
)

func TestUserService_validateUserStruct(t *testing.T) {
	type fields struct {
		storage  *storage.Storage
		l        *zap.Logger
		validate *validator.Validate
	}
	type args struct {
		user model.User
	}

	field := fields{nil, zap.NewExample(), validator.New()}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"success", field, args{model.User{FirstName: "Aybek", LastName: "Zhaisanbay"}}, false},
		{"failed: name contains numbers", field, args{model.User{FirstName: "At78", LastName: "zhayt"}}, true},
		{"failed: FirstName field empty", field, args{model.User{FirstName: "", LastName: "zhayt"}}, true},
		{"failed: FirstName too long", field,
			args{model.User{FirstName: "qwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnmq", LastName: "zhayt"}}, true},
		{"failed: LastName too short", field,
			args{model.User{FirstName: "qwert", LastName: "zh"}}, true},
		{"failed: FirstName has space", field, args{model.User{FirstName: "Aybek", LastName: "Zhaisanbay"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserService{
				storage:  tt.fields.storage,
				l:        tt.fields.l,
				validate: tt.fields.validate,
			}
			if err := s.validateUserStruct(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("validateUserStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
