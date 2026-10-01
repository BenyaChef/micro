package userrepomodel

import (
	"time"

	userentity "github.com/BenyaChef/micro/auth/domain/entity/user"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	emailprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/email"
	timeprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/time"
)

const Columns = "id, email, password_hash, created_at, updated_at"

type User struct {
	ID           string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (m *User) ScanTargets() []any {
	return []any{&m.ID, &m.Email, &m.PasswordHash, &m.CreatedAt, &m.UpdatedAt}
}

func (m *User) InsertValues() []any {
	return []any{m.ID, m.Email, m.PasswordHash, m.CreatedAt, m.UpdatedAt}
}

func ToModel(user *userentity.User) *User {
	return &User{
		ID:           user.ID().String(),
		Email:        user.Email().String(),
		PasswordHash: user.PasswordHash().String(),
		CreatedAt:    user.CreatedAt().Time(),
		UpdatedAt:    user.UpdatedAt().Time(),
	}
}

func ToEntity(model *User) (*userentity.User, error) {
	userID, err := commonuserentity.UserIDFrom(model.ID)
	if err != nil {
		return nil, err
	}

	email, err := emailprimitive.EmailFrom(model.Email)
	if err != nil {
		return nil, err
	}

	passwordHash, err := userentity.PasswordHashFrom(model.PasswordHash)
	if err != nil {
		return nil, err
	}

	return userentity.NewBuilder().
		ID(userID).
		Email(email).
		PasswordHash(passwordHash).
		CreatedAt(timeprimitive.TimestampFrom(model.CreatedAt)).
		UpdatedAt(timeprimitive.TimestampFrom(model.UpdatedAt)).
		Build()
}
