package service

import (
	"go.uber.org/zap"
	"reflect"
	"regexp"
	"testing"
)

func TestAnalysisService_FindEmails(t *testing.T) {
	type fields struct {
		l       *zap.Logger
		emailRx *regexp.Regexp
	}
	type args struct {
		input string
	}
	field := fields{
		zap.NewExample(),
		regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`),
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{"success", field, args{input: "Email:                    \n\n\n    example@gmail.org    email: exe@gmail.com" +
			"ada weew Email: exampletwo@mailru  Email: try@mail.ru"}, []string{"example@gmail.org", "try@mail.ru"}},
		{"success", field, args{input: "Email:                    \n\n\n    example@gmailorg    email: exe@gmail.com" +
			"ada weew Email: exampletwo@mailru  Emaill: try@mail.ru"}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AnalysisService{
				l:       tt.fields.l,
				emailRx: tt.fields.emailRx,
			}
			if got := s.FindEmails(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindEmails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validIIN(t *testing.T) {
	type args struct {
		iin string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"success", args{iin: "011127550738"}, true},
		{"failed", args{iin: "451312550789"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validIIN(tt.args.iin); got != tt.want {
				t.Errorf("validIIN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnalysisService_FindsIINs(t *testing.T) {
	type fields struct {
		l       *zap.Logger
		emailRx *regexp.Regexp
	}
	type args struct {
		input string
	}

	field := fields{zap.NewNop(), nil}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{"success", field, args{input: "asdfewi ewvhvbei evhweuib       \n\n\n  011127550738 ecwe 020304550896"}, []string{"011127550738"}},
		{"success", field, args{input: "asdfewi ewvhvbei evhweuib       \n\n\n  011127550738, ecwe 020304550896"}, []string{"011127550738"}},
		{"success", field, args{input: "asdfewi ewvhvbei evhweuib       \n\n\n  011327550738 ecwe 020304550896"}, []string{}},
		{"success", field, args{input: "asdfewi ewvhvbei evhweuib       \n\n\n wfcwq wqw"}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AnalysisService{
				l:       tt.fields.l,
				emailRx: tt.fields.emailRx,
			}
			if got := s.FindsIINs(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindsIINs() = %v, want %v", got, tt.want)
			}
		})
	}
}
