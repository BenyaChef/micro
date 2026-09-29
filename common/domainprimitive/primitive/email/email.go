package emailprimitive

import (
	"net/mail"
	"strings"
)

type Email string

func EmailFrom(rawEmail string) (Email, error) {
	trimmed := strings.TrimSpace(rawEmail)
	if trimmed == "" {
		return "", ErrEmailIsEmpty
	}
	address, err := mail.ParseAddress(trimmed)
	if err != nil || address.Name != "" {
		return "", ErrEmailInvalidFormat
	}

	return Email(strings.ToLower(address.Address)), nil
}

func (e Email) String() string {
	return string(e)
}
