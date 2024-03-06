package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/zhayt/cert-tz/internal/model"
	"testing"
)

func TestValidateService_validateStruct(t *testing.T) {
	type fields struct {
		validate *validator.Validate
	}
	type args struct {
		data interface{}
	}

	field := fields{validate: validator.New()}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"success model.User", field, args{model.User{FirstName: "Aybek", LastName: "Zhaisanbay"}}, false},
		{"failed model.User: name contains numbers", field, args{model.User{FirstName: "At78", LastName: "zhayt"}}, true},
		{"failed model.User: FirstName field empty", field, args{model.User{FirstName: "", LastName: "zhayt"}}, true},
		{"failed model.User: FirstName too long", field,
			args{model.User{FirstName: "qwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnmq", LastName: "zhayt"}}, true},
		{"failed model.User: LastName too short", field,
			args{model.User{FirstName: "qwert", LastName: "zh"}}, true},
		{"failed model.User: FirstName has space", field, args{model.User{FirstName: "Ayb ek", LastName: "Zhaisanbay"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ValidateService{
				validate: tt.fields.validate,
			}
			if err := s.validateStruct(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("validateStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateService_validateVariable(t *testing.T) {
	type fields struct {
		validate *validator.Validate
	}
	type args struct {
		data interface{}
		tag  string
	}

	field := fields{validate: validator.New()}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"success tag <alphanum>", field, args{"asdQWE145", "alphanum"}, false},
		{"failed tag <alphanum>: contain space", field, args{"asd asd1", "alphanum"}, true},
		{"failed tag <alphanum>: contain symbols", field, args{"asdasd1@", "alphanum"}, true},
		{"failed tag <alphanum>: contain cyrillic alphabet letter", field, args{"asdasd1@ЦАРКА", "alphanum"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ValidateService{
				validate: tt.fields.validate,
			}
			if err := s.validateVariable(tt.args.data, tt.args.tag); (err != nil) != tt.wantErr {
				t.Errorf("validateVariable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
