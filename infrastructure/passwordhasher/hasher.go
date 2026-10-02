package passwordhasher

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	passwordhasherinterface "github.com/BenyaChef/micro/infrastructure/passwordhasher/interface"
)

var _ passwordhasherinterface.PasswordHasher = (*Hasher)(nil)

type Hasher struct {
	cost int
}

func (h *Hasher) Hash(rawPassword string) (string, error) {
	if rawPassword == "" {
		return "", ErrPasswordIsEmpty
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), h.cost)
	if err != nil {
		return "", ErrPasswordHashFailed(err)
	}

	return string(passwordHash), nil
}

func (h *Hasher) Compare(passwordHash, rawPassword string) error {
	if rawPassword == "" {
		return ErrPasswordIsEmpty
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(rawPassword)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrPasswordMismatch
		}

		return ErrPasswordHashFailed(err)
	}

	return nil
}
