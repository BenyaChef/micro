package postentity

import (
	"errors"

	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	timeprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/time"
)

type Builder struct {
	id        PostID
	authorID  commonuserentity.UserID
	title     Title
	body      Body
	createdAt timeprimitive.Timestamp
	updatedAt timeprimitive.Timestamp
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) ID(id PostID) *Builder {
	b.id = id

	return b
}

func (b *Builder) AuthorID(authorID commonuserentity.UserID) *Builder {
	b.authorID = authorID

	return b
}

func (b *Builder) Title(title Title) *Builder {
	b.title = title

	return b
}

func (b *Builder) Body(body Body) *Builder {
	b.body = body

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

func (b *Builder) Build() (*Post, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	b.fillDefaultFields()

	return b.createFromBuilder(), nil
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.id.IsZero() {
		errs = append(errs, ErrPostIDIsRequired)
	}

	if b.authorID.IsZero() {
		errs = append(errs, ErrAuthorIDIsRequired)
	}

	if b.title.IsZero() {
		errs = append(errs, ErrTitleIsRequired)
	}

	if b.body.IsZero() {
		errs = append(errs, ErrBodyIsRequired)
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

func (b *Builder) createFromBuilder() *Post {
	return &Post{
		id:        b.id,
		authorID:  b.authorID,
		title:     b.title,
		body:      b.body,
		createdAt: b.createdAt,
		updatedAt: b.updatedAt,
	}
}
