package userrepomodel

import (
	"time"

	userentity "github.com/BenyaChef/micro/auth/domain/entity/user"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	emailprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/email"
	timeprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/time"
)

type User struct {
	ID           string    `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
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
