package jwtinterface

import "time"

type JWTService interface {
	Issue(subject string) (string, time.Duration, error)
	Subject(rawToken string) (string, error)
}
