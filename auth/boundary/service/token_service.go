package serviceinterface

import "time"

type TokenService interface {
	Issue(subject string) (rawToken string, expiresIn time.Duration, err error)
	Subject(rawToken string) (string, error)
}
