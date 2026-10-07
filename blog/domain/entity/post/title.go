package postentity

import "strings"

const MaxTitleLength = 200

type Title string

func TitleFrom(value string) (Title, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", ErrTitleIsEmpty
	}

	if len([]rune(trimmed)) > MaxTitleLength {
		return "", ErrTitleTooLong(len([]rune(trimmed)), MaxTitleLength)
	}

	return Title(trimmed), nil
}

func (t Title) String() string {
	return string(t)
}

func (t Title) IsZero() bool {
	return t == ""
}
