package postentity

import "strings"

const MaxBodyLength = 20000

type Body string

func BodyFrom(value string) (Body, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", ErrBodyIsEmpty
	}

	if len([]rune(trimmed)) > MaxBodyLength {
		return "", ErrBodyTooLong(len([]rune(trimmed)), MaxBodyLength)
	}

	return Body(trimmed), nil
}

func (b Body) String() string {
	return string(b)
}

func (b Body) IsZero() bool {
	return b == ""
}
