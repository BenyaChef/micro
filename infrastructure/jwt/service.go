package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	jwtinterface "github.com/BenyaChef/micro/infrastructure/jwt/interface"
)

var _ jwtinterface.JWTService = (*Service)(nil)

type Service struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func (s *Service) Issue(subject string) (string, time.Duration, error) {
	if subject == "" {
		return "", 0, ErrSubjectIsEmpty
	}

	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   subject,
		Issuer:    s.issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
	}

	rawToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
	if err != nil {
		return "", 0, ErrSignFailed(err)
	}

	return rawToken, s.ttl, nil
}

func (s *Service) Subject(rawToken string) (string, error) {
	claims := &jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(rawToken, claims, s.keyFunc,
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(s.issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", ErrTokenExpired
		}

		return "", ErrTokenInvalid
	}

	if claims.Subject == "" {
		return "", ErrSubjectIsEmpty
	}

	return claims.Subject, nil
}

func (s *Service) keyFunc(_ *jwt.Token) (any, error) {
	return s.secret, nil
}
