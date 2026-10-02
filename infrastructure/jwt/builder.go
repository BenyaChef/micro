package jwt

import (
	"errors"
	"time"

	jwtmodel "github.com/BenyaChef/micro/infrastructure/jwt/model"
)

type Builder struct {
	secret string
	issuer string
	ttl    time.Duration
}

func NewBuilder() *Builder {
	return &Builder{ttl: jwtmodel.DefaultTTL}
}

func (b *Builder) Secret(secret string) *Builder {
	b.secret = secret

	return b
}

func (b *Builder) Issuer(issuer string) *Builder {
	b.issuer = issuer

	return b
}

func (b *Builder) TTL(ttl time.Duration) *Builder {
	b.ttl = ttl

	return b
}

func (b *Builder) Build() (*Service, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	if len(b.secret) < jwtmodel.MinSecretLength {
		return nil, ErrSecretTooShort(len(b.secret), jwtmodel.MinSecretLength)
	}

	return &Service{secret: []byte(b.secret), issuer: b.issuer, ttl: b.ttl}, nil
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.secret == "" {
		errs = append(errs, ErrSecretIsRequired)
	}

	if b.issuer == "" {
		errs = append(errs, ErrIssuerIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
