package service

import (
	"go.uber.org/zap"
	"regexp"
)

var EmailRX = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

//go:generate mockery --name IAnalysisService
type IAnalysisService interface {
	FindEmails(input string) []string
	FindsIINs(input string) []string
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

// Расшифровка ИИН :
// первые 6 разрядов - это дата рождения ггммдд, то есть 12 августа 1985 года в ИИНе будет 850812
// 7 разряд отвечает за век рождения и пол. Если цифра нечетная - пол мужской, четная - женский. 1,2 - девятнадцатый век, 3,4 - двадцатый, 5,6 - двадцать первый.
// 8-11 разряды - заполняет орган Юстиции.
// 12 разряд - контрольная цифра, которая расчитывается по определенному алгоритму https://titus.kz/?previd=23712

func (s *AnalysisService) FindsIINs(input string) []string {
	// делаем грубый отбор, отбор по длине и по "коректной дате" месяц не больше 12, день не больше 31
	regex := regexp.MustCompile(`\b\d{2}((0[1-9]|1[0-2])(0[1-9]|1[0-9]|2[0-9]|30|31)([1-6]))\d{5}\b`)
	matches := regex.FindAllString(input, -1)

	iins := make([]string, 0, len(matches))
	for _, match := range matches {
		// проверка на валидность иин-а
		if validIIN(match) {
			iins = append(iins, match)
		}
	}

	return iins
}

// проверка валидности иин-а
// иин считается валидным когда сумма произведения порядка разряда на его значение по mod 11
// равно контрольная цифре. Контрольная цыфра, это стоящии в 12 разряде
//
// Пример: 011127550738
//
// Находим значения сумму произведения порядка разряда на его значение по mod 11
// 1 * 0 + 2 * 1 + 3 * 1 + 4 * 1 + 5 * 2 + 6 * 7 + 7 * 5 + 8 * 5 + 9 * 0 + 10 * 7 + 11 * 3
// 0 + 2 + 3 + 4 + 10 + 42 + 35 + 40 + 0 + 70 + 33 = 239
// 239mod11 = 8
//
// Сравним его с контрольной цифрой который, находится в 12 разряде
//
// 8 == 8
// Ответ: ИИН валидный

func validIIN(iin string) bool {
	digits := []byte(iin)
	sum := 0

	// сумма произведения порядка разряда на его значение до контрольной цифры
	for i := 0; i < 11; i++ {
		sum += (i + 1) * int(digits[i]-'0')
	}

	controlDigit := int(digits[11] - '0')
	mod := sum % 11
	// проверяем остаток от деления от 11 с контрольной цифрой
	// если они равный иин коректный
	if mod == controlDigit {
		return true
	}

	return false
}
