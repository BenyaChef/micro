package tokenusecasemock

import (
	"time"

	serviceinterface "github.com/BenyaChef/micro/auth/boundary/service"
)

var _ serviceinterface.TokenService = (*TokenServiceMock)(nil)

const (
	StubRawToken  = "stub-raw-token"
	StubExpiresIn = time.Hour
)

type TokenServiceMock struct {
	IssueFunc   func(subject string) (string, time.Duration, error)
	SubjectFunc func(rawToken string) (string, error)
}

func NewTokenServiceMock() *TokenServiceMock {
	return &TokenServiceMock{
		IssueFunc: func(string) (string, time.Duration, error) {
			return StubRawToken, StubExpiresIn, nil
		},
		SubjectFunc: func(string) (string, error) {
			return "550e8400-e29b-41d4-a716-446655440000", nil
		},
	}
}

func (m *TokenServiceMock) Issue(subject string) (string, time.Duration, error) {
	return m.IssueFunc(subject)
}

func (m *TokenServiceMock) Subject(rawToken string) (string, error) {
	return m.SubjectFunc(rawToken)
}
