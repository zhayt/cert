package service

import (
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestSubStrService_LongestSubstring(t *testing.T) {
	type fields struct {
		validate *ValidateService
		l        *zap.Logger
	}
	type args struct {
		str string
	}

	field := fields{
		validate: NewValidateService(),
		l:        zap.NewExample(),
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{"success", field, args{"pwwkew"}, "wke", false},
		{"success", field, args{"zzhaytQWER123zhayt"}, "zhaytQWER123", false},
		{"field: invalid data", field, args{" "}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SubStrService{
				validate: tt.fields.validate,
				l:        tt.fields.l,
			}
			got, err := s.LongestSubstring(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("LongestSubstring() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LongestSubstring() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cleanHash(t *testing.T) {
	type args struct {
		s      string
		letter map[byte]int
	}
	tests := []struct {
		name    string
		args    args
		wantMap map[byte]int
	}{
		{"success", args{s: "asd", letter: map[byte]int{'a': 0, 's': 1, 'd': 2, 'e': 3}}, map[byte]int{'e': 3}},
		{"success", args{s: "asd", letter: map[byte]int{'a': 0, 's': 1, 'd': 2}}, map[byte]int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanHash(tt.args.s, tt.args.letter)
			if !reflect.DeepEqual(tt.wantMap, tt.args.letter) {
				t.Errorf("cleanHash() map = %v, wantMap %v", tt.args.letter, tt.wantMap)
				return
			}
		})
	}
}

func Test_longestString(t *testing.T) {
	type args struct {
		s1 string
		s2 string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"success", args{s1: "asd", s2: "asdf"}, "asdf"},
		{"success", args{s1: "asdqwe", s2: "asdf"}, "asdqwe"},
		{"success", args{s1: "", s2: "asdf"}, "asdf"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := longestString(tt.args.s1, tt.args.s2); got != tt.want {
				t.Errorf("longestString() = %v, want %v", got, tt.want)
			}
		})
	}
}
