package userentity

import (
	"errors"

	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	emailprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/email"
	timeprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/time"
)

type Builder struct {
	id           commonuserentity.UserID
	email        emailprimitive.Email
	passwordHash PasswordHash
	createdAt    timeprimitive.Timestamp
	updatedAt    timeprimitive.Timestamp
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) ID(id commonuserentity.UserID) *Builder {
	b.id = id

	return b
}

func (b *Builder) Email(email emailprimitive.Email) *Builder {
	b.email = email

	return b
}

func (b *Builder) PasswordHash(passwordHash PasswordHash) *Builder {
	b.passwordHash = passwordHash

	return b
}

func (b *Builder) CreatedAt(createdAt timeprimitive.Timestamp) *Builder {
	b.createdAt = createdAt

	return b
}

func (b *Builder) UpdatedAt(updatedAt timeprimitive.Timestamp) *Builder {
	b.updatedAt = updatedAt

	return b
}

func (b *Builder) Build() (*User, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	b.fillDefaultFields()

	return b.createFromBuilder(), nil
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.id.IsZero() {
		errs = append(errs, ErrUserIDIsRequired)
	}

	if b.email == "" {
		errs = append(errs, ErrEmailIsRequired)
	}

	if b.passwordHash.IsZero() {
		errs = append(errs, ErrPasswordHashIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (b *Builder) fillDefaultFields() {
	now := timeprimitive.Now()

	if b.createdAt.IsZero() {
		b.createdAt = now
	}

	if b.updatedAt.IsZero() {
		b.updatedAt = now
	}
}

func (b *Builder) createFromBuilder() *User {
	return &User{
		id:           b.id,
		email:        b.email,
		passwordHash: b.passwordHash,
		createdAt:    b.createdAt,
		updatedAt:    b.updatedAt,
	}
}
