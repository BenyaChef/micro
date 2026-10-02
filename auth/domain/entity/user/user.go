package userentity

import (
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	emailprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/email"
	timeprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/time"
)

type User struct {
	id           commonuserentity.UserID
	email        emailprimitive.Email
	passwordHash PasswordHash
	createdAt    timeprimitive.Timestamp
	updatedAt    timeprimitive.Timestamp
}

func (u *User) ID() commonuserentity.UserID {
	return u.id
}

func (u *User) Email() emailprimitive.Email {
	return u.email
}

func (u *User) PasswordHash() PasswordHash {
	return u.passwordHash
}

func (u *User) CreatedAt() timeprimitive.Timestamp {
	return u.createdAt
}

func (u *User) UpdatedAt() timeprimitive.Timestamp {
	return u.updatedAt
}
