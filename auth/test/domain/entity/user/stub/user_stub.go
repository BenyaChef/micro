package userentitystub

import (
	userentity "github.com/BenyaChef/micro/auth/domain/entity/user"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	emailprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/email"
)

const (
	RawEmail        = "user@example.com"
	RawPasswordHash = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
)

func GetEmail() emailprimitive.Email {
	email, err := emailprimitive.EmailFrom(RawEmail)
	if err != nil {
		panic(err)
	}

	return email
}

func GetPasswordHash() userentity.PasswordHash {
	passwordHash, err := userentity.PasswordHashFrom(RawPasswordHash)
	if err != nil {
		panic(err)
	}

	return passwordHash
}

func GetUserID() commonuserentity.UserID {
	return commonuserentity.NewUserID()
}

func GetUserBuilder() *userentity.Builder {
	return userentity.NewBuilder().
		ID(GetUserID()).
		Email(GetEmail()).
		PasswordHash(GetPasswordHash())
}

func GetValidUser() *userentity.User {
	user, err := GetUserBuilder().Build()
	if err != nil {
		panic(err)
	}

	return user
}
