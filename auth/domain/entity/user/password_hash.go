package userentity

type PasswordHash string

func PasswordHashFrom(value string) (PasswordHash, error) {
	if value == "" {
		return "", ErrPasswordHashIsEmpty
	}

	return PasswordHash(value), nil
}

func (h PasswordHash) String() string {
	return string(h)
}

func (h PasswordHash) IsZero() bool {
	return h == ""
}
