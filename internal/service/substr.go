package service

import (
	"go.uber.org/zap"
)

//go:generate mockery --name ISubStrService
type ISubStrService interface {
	LongestSubstring(str string) (string, error)
}

type SubStrService struct {
	validate *ValidateService
	l        *zap.Logger
}

func (s *SubStrService) LongestSubstring(str string) (string, error) {
	if err := s.validate.validateVariable(str, "alphanum"); err != nil {
		s.l.Error("validateVariable error", zap.Error(err))
		return "", ErrInvalidData
	}

	var res string
	var i, j int
	letter := make(map[byte]int)

	for ; j < len(str); j++ {
		if end, ok := letter[str[j]]; ok {
			res = longestString(res, str[i:j])

			cleanHash(str[i:end+1], letter)
			i = end + 1
		}
		letter[str[j]] = j
	}

	return longestString(res, str[i:j]), nil
}

func cleanHash(s string, letter map[byte]int) {
	for i := 0; i < len(s); i++ {
		delete(letter, s[i])
	}
}

func longestString(s1, s2 string) string {
	if len(s1) >= len(s2) {
		return s1
	}

	return s2
}

func NewSubStrService(validate *ValidateService, l *zap.Logger) *SubStrService {
	return &SubStrService{validate: validate, l: l}
}
