package userusecasemock

import serviceinterface "github.com/BenyaChef/micro/auth/boundary/service"

var _ serviceinterface.PasswordHasher = (*PasswordHasherMock)(nil)

const StubPasswordHash = "hashed"

type PasswordHasherMock struct {
	HashFunc    func(rawPassword string) (string, error)
	CompareFunc func(passwordHash, rawPassword string) error
}

func NewPasswordHasherMock() *PasswordHasherMock {
	return &PasswordHasherMock{
		HashFunc:    func(string) (string, error) { return StubPasswordHash, nil },
		CompareFunc: func(string, string) error { return nil },
	}
}

func (m *PasswordHasherMock) Hash(rawPassword string) (string, error) {
	return m.HashFunc(rawPassword)
}

func (m *PasswordHasherMock) Compare(passwordHash, rawPassword string) error {
	return m.CompareFunc(passwordHash, rawPassword)
}
