package service

import (
	"go.uber.org/zap"
	"regexp"
)

var EmailRX = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

//go:generate mockery --name IAnalysisService
type IAnalysisService interface {
	FindEmails(input string) []string
}

type AnalysisService struct {
	l       *zap.Logger
	emailRx *regexp.Regexp
}

func NewAnalysisService(l *zap.Logger) *AnalysisService {
	emailRx := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	return &AnalysisService{
		l:       l,
		emailRx: emailRx,
	}
}

func (s *AnalysisService) FindEmails(input string) []string {
	regex := regexp.MustCompile(`Email:\s+(\S+)`)
	matches := regex.FindAllStringSubmatch(input, -1)

	emails := make([]string, 0, len(matches))
	for _, match := range matches {
		if s.emailRx.MatchString(match[1]) {
			emails = append(emails, match[1])
		}
	}

	return emails
}
